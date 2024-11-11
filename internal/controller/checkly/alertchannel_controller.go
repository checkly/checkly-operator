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

	corev1 "k8s.io/api/core/v1"
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

// AlertChannelReconciler reconciles a AlertChannel object
type AlertChannelReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	ApiClient        checkly.Client
	ControllerDomain string
}

//+kubebuilder:rbac:groups=k8s.checklyhq.com,resources=alertchannels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.checklyhq.com,resources=alertchannels/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.checklyhq.com,resources=alertchannels/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *AlertChannelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciler started")

	acFinalizer := fmt.Sprintf("%s/finalizer", r.ControllerDomain)

	ac := &checklyv1alpha1.AlertChannel{}

	err := r.Get(ctx, req.NamespacedName, ac)

	// ////////////////////////////////
	// Delete Logic
	// ///////////////////////////////
	if err != nil {
		if errors.IsNotFound(err) {
			// The resource has been deleted
			logger.Info("Deleted", "checkly AlertChannel ID", ac.Status.ID)
			return ctrl.Result{}, nil
		}
		// Error reading the object
		logger.Error(err, "can't read the object")
		return ctrl.Result{}, nil
	}

	// ////////////////////////////////
	// Remove Finalizer Logic
	// ///////////////////////////////

	if ac.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(ac, acFinalizer) {
			logger.Info("Finalizer is present, trying to delete Checkly AlertChannel", "ID", ac.Status.ID)
			err := external.DeleteAlertChannel(ac, r.ApiClient)
			if err != nil {
				logger.Error(err, "Failed to delete checkly AlertChannel")
				return ctrl.Result{}, err
			}

			logger.Info("Successfully deleted checkly AlertChannel", "ID", ac.Status.ID)

			controllerutil.RemoveFinalizer(ac, acFinalizer)
			err = r.Update(ctx, ac)
			if err != nil {
				return ctrl.Result{}, err
			}
			logger.Info("Successfully deleted finalizer from AlertChannel")
		}
		return ctrl.Result{}, nil
	}

	// /////////////////////////////
	// Add Finalizer logic
	// ////////////////////////////
	if !controllerutil.ContainsFinalizer(ac, acFinalizer) {
		controllerutil.AddFinalizer(ac, acFinalizer)
		err = r.Update(ctx, ac)
		if err != nil {
			logger.Error(err, "Failed to update AlertChannel status")
			return ctrl.Result{}, err
		}
		logger.Info("Added finalizer", "checkly AlertChannel ID", ac.Status.ID)
		return ctrl.Result{}, nil
	}

	// /////////////////////////////
	// OpsGenie logic + secret retrieval
	// ////////////////////////////
	opsGenieConfig := checkly.AlertChannelOpsgenie{}
	if ac.Spec.OpsGenie.APISecret != (corev1.ObjectReference{}) {
		secret := &corev1.Secret{}
		err := r.Get(ctx,
			types.NamespacedName{
				Name:      ac.Spec.OpsGenie.APISecret.Name,
				Namespace: ac.Spec.OpsGenie.APISecret.Namespace},
			secret)
		if err != nil {
			logger.Info("Unable to read secret for API Key", "err", err)
			return ctrl.Result{}, err
		}

		secretValue := string(secret.Data[ac.Spec.OpsGenie.APISecret.FieldPath])
		if secretValue == "" {
			logger.Info("Secret value is empty")
			return ctrl.Result{}, err
		}

		opsGenieConfig = checkly.AlertChannelOpsgenie{
			Name:     ac.Name,
			APIKey:   secretValue,
			Region:   ac.Spec.OpsGenie.Region,
			Priority: ac.Spec.OpsGenie.Priority,
		}

	}

	// /////////////////////////////
	// Update logic
	// ////////////////////////////

	// Determine if it's a new object or if it's an update to an existing object
	if ac.Status.ID != 0 {
		// Existing object, we need to update it
		logger.Info("Existing object, with ID", "checkly AlertChannel ID", ac.Status.ID)
		err := external.UpdateAlertChannel(ac, opsGenieConfig, r.ApiClient)
		if err != nil {
			logger.Error(err, "Failed to update checkly AlertChannel")
			return ctrl.Result{}, err
		}
		logger.Info("Updated checkly AlertChannel", "ID", ac.Status.ID)
		return ctrl.Result{}, nil
	}

	// /////////////////////////////
	// Create logic
	// ////////////////////////////
	acID, err := external.CreateAlertChannel(ac, opsGenieConfig, r.ApiClient)
	if err != nil {
		logger.Error(err, "Failed to create checkly AlertChannel")
		return ctrl.Result{}, err
	}

	// Update the custom resource Status with the returned ID
	ac.Status.ID = acID
	err = r.Status().Update(ctx, ac)
	if err != nil {
		logger.Error(err, "Failed to update AlertChannel status", "ID", ac.Status.ID)
		return ctrl.Result{}, err
	}
	logger.Info("New checkly AlertChannel created", "ID", ac.Status.ID)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AlertChannelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&checklyv1alpha1.AlertChannel{}).
		Complete(r)
}
