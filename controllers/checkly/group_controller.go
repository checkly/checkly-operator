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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	checklyv1alpha1 "github.com/imgarena/checkly-operator/apis/checkly/v1alpha1"
	external "github.com/imgarena/checkly-operator/external/checkly"
)

// GroupReconciler reconciles a Group object
type GroupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=checkly.imgarena.com,resources=groups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=checkly.imgarena.com,resources=groups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=checkly.imgarena.com,resources=groups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Group object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *GroupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciler started")

	groupFinalizer := "checkly.imgarena.com/finalizer"

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
		logger.Error(err, "can't read the object")
		return ctrl.Result{}, nil
	}

	if group.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(group, groupFinalizer) {
			logger.Info("Finalizer is present, trying to delete Checkly group", "checkly ID", group.Status.ID)
			err := external.GroupDelete(group.Status.ID)
			if err != nil {
				logger.Error(err, "Failed to delete checkly group")
				return ctrl.Result{}, err
			}

			logger.Info("Successfully deleted checkly group", "checkly ID", group.Status.ID)

			controllerutil.RemoveFinalizer(group, groupFinalizer)
			err = r.Update(ctx, group)
			if err != nil {
				return ctrl.Result{}, err
			}
			logger.Info("Successfully deleted finalizer")
		}
		return ctrl.Result{}, nil
	}

	// Object found, let's do something with it. It's either updated, or it's new.
	logger.Info("Object found")

	// /////////////////////////////
	// Finalizer logic
	// ////////////////////////////
	if !controllerutil.ContainsFinalizer(group, groupFinalizer) {
		controllerutil.AddFinalizer(group, groupFinalizer)
		err = r.Update(ctx, group)
		if err != nil {
			logger.Error(err, "Failed to update group status")
			return ctrl.Result{}, err
		}
		logger.Info("Added finalizer", "checkly ID", group.Status.ID)
		return ctrl.Result{}, nil
	}

	// Create internal Check type
	internalCheck := external.Group{
		Name:          group.Name,
		Namespace:     group.Namespace,
		Activated:     group.Spec.Activated,
		Locations:     group.Spec.Locations,
		AlertChannels: group.Spec.AlertChannels,
		ID:            group.Status.ID,
	}

	// /////////////////////////////
	// Update logic
	// ////////////////////////////

	// Determine if it's a new object or if it's an update to an existing object
	if group.Status.ID != 0 {
		// Existing object, we need to update it
		logger.Info("Existing object, with ID", "checkly ID", group.Status.ID)
		err := external.GroupUpdate(internalCheck)
		// err :=
		if err != nil {
			logger.Error(err, "Failed to update the checkly group")
			return ctrl.Result{}, err
		}
		logger.Info("Updated checkly check", "checkly ID", group.Status.ID)
		return ctrl.Result{}, nil
	}

	// /////////////////////////////
	// Create logic
	// ////////////////////////////

	checklyID, err := external.GroupCreate(internalCheck)
	if err != nil {
		logger.Error(err, "Failed to create checkly group")
		return ctrl.Result{}, nil
	}

	// Update the custom resource Status with the returned ID

	group.Status.ID = checklyID
	err = r.Status().Update(ctx, group)
	if err != nil {
		logger.Error(err, "Failed to update group status")
		return ctrl.Result{}, err
	}
	logger.Info("New checkly group created", "ID", group.Status.ID)

	return ctrl.Result{}, nil
	// TODO(user): your logic here

	// return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&checklyv1alpha1.Group{}).
		Complete(r)
}