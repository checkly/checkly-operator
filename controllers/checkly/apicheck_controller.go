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
	"sigs.k8s.io/controller-runtime/pkg/log"

	checklyv1alpha1 "github.com/imgarena/checkly-operator/apis/checkly/v1alpha1"
)

// ApiCheckReconciler reconciles a ApiCheck object
type ApiCheckReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=checkly.imgarena.com,resources=apichecks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=checkly.imgarena.com,resources=apichecks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=checkly.imgarena.com,resources=apichecks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ApiCheck object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ApiCheckReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	apiCheck := &checklyv1alpha1.ApiCheck{}
	err := r.Get(ctx, req.NamespacedName, apiCheck)

	if err != nil {
		if errors.IsNotFound(err) {
			// Probably the resource has been deleted
			log.Log.Info("Deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object
		return ctrl.Result{}, nil
	}

	// Object found, let's do something with it.
	log.Log.Info("Object found", "endpoint", apiCheck.Spec.Endpoint)

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApiCheckReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&checklyv1alpha1.ApiCheck{}).
		Complete(r)
}
