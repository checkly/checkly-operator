/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package checkly

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/checkly/checkly-go-sdk"
	checklyv1alpha1 "github.com/checkly/checkly-operator/api/checkly/v1alpha1"
	external "github.com/checkly/checkly-operator/external/checkly"
)

// GroupReconciler reconciles a Group object
type GroupReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	ApiClient        checkly.Client
	ControllerDomain string
}

//+kubebuilder:rbac:groups=k8s.checklyhq.com,resources=groups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.checklyhq.com,resources=groups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.checklyhq.com,resources=groups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *GroupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.V(1).Info("Reconciler started")

	groupFinalizer := fmt.Sprintf("%s/finalizer", r.ControllerDomain)

	group := &checklyv1alpha1.Group{}

	// ////////////////////////////////
	// Delete Logic
	// TODO: Add logic to determine if there are any checks that are part of the group. If yes, throw error and do not delete the group until the checks have been deleted first.
	// ///////////////////////////////
	err := r.Get(ctx, req.NamespacedName, group)
	if err != nil {
		if errors.IsNotFound(err) {
			// The resource has been deleted
			logger.Info("Deleted", "group ID", group.Status.ID, "name", group.Name)
			return ctrl.Result{}, nil
		}
		// Error reading the object
		logger.Error(err, "can't read the Group object")
		return ctrl.Result{}, nil
	}

	// If DeletionTimestamp is present, the object is marked for deletion, we need to remove the finalizer
	if group.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(group, groupFinalizer) {
			logger.V(1).Info("Finalizer is present, trying to delete Checkly group", "checkly group ID", group.Status.ID)
			err := external.GroupDelete(group.Status.ID, r.ApiClient)
			if err != nil {
				logger.Error(err, "Failed to delete checkly group")
				return ctrl.Result{}, err
			}

			logger.Info("Successfully deleted checkly group", "checkly group ID", group.Status.ID)

			controllerutil.RemoveFinalizer(group, groupFinalizer)
			err = r.Update(ctx, group)
			if err != nil {
				logger.Error(err, "Failed to delete finalizer")
				return ctrl.Result{}, err
			}
			logger.V(1).Info("Successfully deleted finalizer")
		}
		return ctrl.Result{}, nil
	}

	// Object found, let's do something with it. It's either updated, or it's new.
	logger.V(1).Info("Checkly group found")

	// /////////////////////////////
	// Finalizer logic
	// ////////////////////////////
	if !controllerutil.ContainsFinalizer(group, groupFinalizer) {
		controllerutil.AddFinalizer(group, groupFinalizer)
		err = r.Update(ctx, group)
		if err != nil {
			logger.Error(err, "Failed to add Group finalizer")
			return ctrl.Result{}, err
		}
		logger.V(1).Info("Added finalizer", "checkly group ID", group.Status.ID)
		return ctrl.Result{}, nil
	}

	// /////////////////////////////
	// AlertChannelsSubscription logic
	// ////////////////////////////
	var alertChannels []checkly.AlertChannelSubscription

	if len(group.Spec.AlertChannels) != 0 {
		for _, alertChannel := range group.Spec.AlertChannels {
			ac := &checklyv1alpha1.AlertChannel{}
			err := r.Get(ctx, types.NamespacedName{Name: alertChannel}, ac)
			if err != nil {
				logger.Error(err, "Could not find alertChannel resource", "name", alertChannel)
				return ctrl.Result{}, err
			}
			if ac.Status.ID == 0 {
				logger.Info("AlertChannel ID not yet populated, we'll retry")
				return ctrl.Result{Requeue: true}, nil
			}
			alertChannels = append(alertChannels, checkly.AlertChannelSubscription{
				ChannelID: ac.Status.ID,
				Activated: true,
			})
		}
	}

	// Iterate and remove keys ending with "/argo-app"
	newLabels := make(map[string]string)
	for key, value := range group.Labels {
	    if !strings.Contains(key, "/argo-app") {
	        newLabels[key] = value
	    }
	}

	// Create internal Check type
	internalCheck := external.Group{
		Name:          group.Name,
		Activated:     group.Spec.Activated,
		Locations:     group.Spec.Locations,
		AlertChannels: alertChannels,
		ID:            group.Status.ID,
		Labels:        newLabels,
	}

	// /////////////////////////////
	// Update logic
	// ////////////////////////////

	// Determine if it's a new object or if it's an update to an existing object
	if group.Status.ID != 0 {
		// Existing object, we need to update it
		logger.V(1).Info("Existing object, with ID", "checkly group ID", group.Status.ID)
		err := external.GroupUpdate(internalCheck, r.ApiClient)
		if err != nil {
			logger.Error(err, "Failed to update the checkly group")
			return ctrl.Result{}, err
		}
		logger.V(1).Info("Updated checkly check", "checkly group ID", group.Status.ID)
		return ctrl.Result{}, nil
	}

	// /////////////////////////////
	// Create logic
	// ////////////////////////////
	checklyID, err := external.GroupCreate(internalCheck, r.ApiClient)
	if err != nil {
		logger.Error(err, "Failed to create checkly group")
		return ctrl.Result{}, err
	}

	// Update the custom resource Status with the returned ID
	group.Status.ID = checklyID
	err = r.Status().Update(ctx, group)
	if err != nil {
		logger.Error(err, "Failed to update group status", "ID", group.Status.ID)
		return ctrl.Result{}, err
	}
	logger.Info("New checkly group created", "ID", group.Status.ID)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&checklyv1alpha1.Group{}).
		Complete(r)
}
