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

package networking

import (
	"context"
	"fmt"

	checklyv1alpha1 "github.com/imgarena/checkly-operator/apis/checkly/v1alpha1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// IngressReconciler reconciles a Ingress object
type IngressReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/finalizers,verbs=update
//+kubebuilder:rbac:groups=checkly.imgarena.com,resources=apichecks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=checkly.imgarena.com,resources=apichecks/status,verbs=get;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciler started")

	ingress := &networkingv1.Ingress{}
	apiCheck := &checklyv1alpha1.ApiCheck{}

	// Check if ingress object is still present
	err := r.Get(ctx, req.NamespacedName, ingress)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Ingress got deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Can't read the Ingress object")
		return ctrl.Result{}, err
	}
	logger.Info("Ingress Object found")

	// Check if annotation is present on the object
	checklyAnnotation := ingress.Annotations["checkly.imgarena.com/enabled"] == "true"
	if !checklyAnnotation {
		// Annotation may have been removed or updated, we have to determine if we need to delete a previously created ApiCheck resource
		logger.Info("annotation is not present, checking if ApiCheck was created")
		err = r.Get(ctx, req.NamespacedName, apiCheck)
		if err != nil {
			logger.Info("Apicheck not present")
			return ctrl.Result{}, nil
		}
		logger.Info("ApiCheck is present, but we need to delete it")
		err = r.Delete(ctx, apiCheck)
		if err != nil {
			logger.Info("Failed to delete ApiCheck")
			return ctrl.Result{}, err
			// }
		}

		return ctrl.Result{}, nil
	}

	// Gather data for the checkly check
	apiCheckSpec, err := r.gatherApiCheckData(ingress)
	if err != nil {
		logger.Info("unable to gather data for the apiCheck resource")
		return ctrl.Result{}, err
	}

	// Check and see if the ApiCheck has been created before
	err = r.Get(ctx, req.NamespacedName, apiCheck)
	if err == nil {
		logger.Info("apiCheck exists, doing an update")
		// We can reference the exiting apiCheck object that the server returned
		apiCheck.Spec = apiCheckSpec
		err = r.Update(ctx, apiCheck)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Create apiCheck
	// We need to write the k8s spec resources as it is a new object
	newApiCheck := &checklyv1alpha1.ApiCheck{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingress.Name,
			Namespace: ingress.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(ingress, networkingv1.SchemeGroupVersion.WithKind("ingress")),
			},
		},
		Spec: apiCheckSpec,
	}

	err = r.Create(ctx, newApiCheck)
	if err != nil {
		logger.Info("Failed to create ApiCheck", "err", err)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1.Ingress{}).
		Complete(r)
}

func (r *IngressReconciler) gatherApiCheckData(ingress *networkingv1.Ingress) (apiCheckSpec checklyv1alpha1.ApiCheckSpec, err error) {

	annotationHost := "checkly.imgarena.com"
	annotationPath := fmt.Sprintf("%s/path", annotationHost)
	annotationEndpoint := fmt.Sprintf("%s/endpoint", annotationHost)
	annotationSuccess := fmt.Sprintf("%s/success", annotationHost)
	annotationGroup := fmt.Sprintf("%s/group", annotationHost)

	// Construct the endpoint
	path := ""
	if ingress.Annotations[annotationPath] != "" {
		path = ingress.Annotations[annotationPath]
	}

	var host string
	if ingress.Annotations[annotationEndpoint] == "" {
		host = ingress.Spec.Rules[0].Host
	} else {
		host = ingress.Annotations[annotationEndpoint]
	}

	endpoint := fmt.Sprintf("https://%s%s", host, path)

	// Expected success code
	var success string
	if ingress.Annotations[annotationSuccess] != "" {
		success = ingress.Annotations[annotationSuccess]
	} else {
		success = "200"
	}

	// Group
	var group string
	if ingress.Annotations[annotationGroup] != "" {
		group = ingress.Annotations[annotationGroup]
	} else {
		err = fmt.Errorf("could not find a value for the group annotation, can't continue without one")
	}

	apiCheckSpec = checklyv1alpha1.ApiCheckSpec{
		Endpoint: endpoint,
		Group:    group,
		Team:     "group-test",
		Success:  success,
	}

	// Last return
	return
}
