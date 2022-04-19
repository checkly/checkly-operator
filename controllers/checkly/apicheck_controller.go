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
	"net/http"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	checklyv1alpha1 "github.com/imgarena/checkly-operator/apis/checkly/v1alpha1"

	checkly "github.com/checkly/checkly-go-sdk"
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

	apiCheckFinalizer := "checkly.imgarena.com/finalizer"

	apiCheck := &checklyv1alpha1.ApiCheck{}

	// ////////////////////////////////
	// Delete Logic
	// ///////////////////////////////
	err := r.Get(ctx, req.NamespacedName, apiCheck)
	if err != nil {
		if errors.IsNotFound(err) {
			// The resource has been deleted
			log.Log.Info("Deleted", "checkly ID", apiCheck.Status.ID, "endpoint", apiCheck.Spec.Endpoint, "name", apiCheck.Name)
			return ctrl.Result{}, nil
		}
		// Error reading the object
		log.Log.Error(err, "can't read the object")
		return ctrl.Result{}, nil
	}

	if apiCheck.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(apiCheck, apiCheckFinalizer) {
			log.Log.Info("Finalizer is present, trying to delete Checkly check", "checkly ID", apiCheck.Status.ID)
			err := delete(apiCheck.Status.ID)
			if err != nil {
				log.Log.Error(err, "Failed to delete checkly API check")
				return ctrl.Result{}, err
			}

			log.Log.Info("Successfully deleted checkly API check", "checkly ID", apiCheck.Status.ID)

			controllerutil.RemoveFinalizer(apiCheck, apiCheckFinalizer)
			err = r.Update(ctx, apiCheck)
			if err != nil {
				return ctrl.Result{}, err
			}
			log.Log.Info("Successfully deleted finalizer")
		}
		return ctrl.Result{}, nil
	}

	// Object found, let's do something with it. It's either updated, or it's new.
	log.Log.Info("Object found", "endpoint", apiCheck.Spec.Endpoint)

	// /////////////////////////////
	// Finalizer logic
	// ////////////////////////////
	if !controllerutil.ContainsFinalizer(apiCheck, apiCheckFinalizer) {
		controllerutil.AddFinalizer(apiCheck, apiCheckFinalizer)
		err = r.Update(ctx, apiCheck)
		if err != nil {
			log.Log.Error(err, "Failed to update ApiCheck status")
			return ctrl.Result{}, err
		}
		log.Log.Info("Added finalizer", "checkly ID", apiCheck.Status.ID, "endpoint", apiCheck.Spec.Endpoint)
		return ctrl.Result{}, nil
	}
	// /////////////////////////////
	// Update logic
	// ////////////////////////////

	// Determine if it's a new object or if it's an update to an existing object
	if apiCheck.Status.ID != "" {
		// Existing object, we need to update it
		log.Log.Info("Existing object, with ID", "checkly ID", apiCheck.Status.ID, "endpoint", apiCheck.Spec.Endpoint)
		err := update(apiCheck)
		if err != nil {
			log.Log.Error(err, "Failed to update the checkly check")
			return ctrl.Result{}, err
		}
		log.Log.Info("Updated checkly check", "checkly ID", apiCheck.Status.ID)
		return ctrl.Result{}, nil
	}

	// /////////////////////////////
	// Create logic
	// ////////////////////////////

	checklyID, err := create(apiCheck)
	if err != nil {
		log.Log.Error(err, "Failed to create checkly alert")
		return ctrl.Result{}, nil
	}

	// Update the custom resource Status with the returned ID

	apiCheck.Status.ID = checklyID
	err = r.Status().Update(ctx, apiCheck)
	if err != nil {
		log.Log.Error(err, "Failed to update ApiCheck status")
		return ctrl.Result{}, err
	}
	log.Log.Info("New checkly check created with", "checkly ID", apiCheck.Status.ID, "endpoint", apiCheck.Spec.Endpoint)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApiCheckReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&checklyv1alpha1.ApiCheck{}).
		Complete(r)
}

func checklyClient() (client checkly.Client, ctx context.Context, cancel context.CancelFunc) {
	baseUrl := "https://api.checklyhq.com"
	apiKey := os.Getenv("CHECKLY_API_KEY")
	accountId := os.Getenv("CHECKLY_ACCOUNT_ID")
	client = checkly.NewClient(
		baseUrl,
		apiKey,
		nil, //custom http client, defaults to http.DefaultClient
		nil, //io.Writer to output debug messages
	)

	client.SetAccountId(accountId)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	return
}

func create(apiCheck *checklyv1alpha1.ApiCheck) (ID string, err error) {

	alertSettings := checkly.AlertSettings{
		EscalationType: checkly.RunBased,
		RunBasedEscalation: checkly.RunBasedEscalation{
			FailedRunThreshold: 5,
		},
		TimeBasedEscalation: checkly.TimeBasedEscalation{
			MinutesFailingThreshold: 5,
		},
		Reminders: checkly.Reminders{
			Interval: 5,
		},
		SSLCertificates: checkly.SSLCertificates{
			Enabled:        false,
			AlertThreshold: 3,
		},
	}

	check := checkly.Check{
		Name:                 apiCheck.Name,
		Type:                 checkly.TypeAPI,
		Frequency:            5,
		DegradedResponseTime: 5000,
		MaxResponseTime:      15000,
		Activated:            false, // deactivated for development
		Muted:                true,  // muted for development
		ShouldFail:           false,
		DoubleCheck:          false,
		SSLCheck:             false,
		LocalSetupScript:     "",
		LocalTearDownScript:  "",
		Locations: []string{
			"eu-west-1",
			"ap-northeast-2",
		},
		Tags: []string{
			"foo",
			"bar",
		},
		AlertSettings:          alertSettings,
		UseGlobalAlertSettings: false,
		Request: checkly.Request{
			Method:  http.MethodGet,
			URL:     apiCheck.Spec.Endpoint,
			Headers: []checkly.KeyValue{
				// {
				// 	Key:   "X-Test",
				// 	Value: "foo",
				// },
			},
			QueryParameters: []checkly.KeyValue{
				// {
				// 	Key:   "query",
				// 	Value: "foo",
				// },
			},
			Assertions: []checkly.Assertion{
				{
					Source:     checkly.StatusCode,
					Comparison: checkly.Equals,
					Target:     apiCheck.Spec.Success,
				},
			},
			Body:     "",
			BodyType: "NONE",
		},
	}

	client, ctx, cancel := checklyClient()
	defer cancel()

	gotCheck, err := client.Create(ctx, check)
	if err != nil {
		return
	}

	ID = gotCheck.ID

	return
}

func update(apiCheck *checklyv1alpha1.ApiCheck) (err error) {

	alertSettings := checkly.AlertSettings{
		EscalationType: checkly.RunBased,
		RunBasedEscalation: checkly.RunBasedEscalation{
			FailedRunThreshold: 5,
		},
		TimeBasedEscalation: checkly.TimeBasedEscalation{
			MinutesFailingThreshold: 5,
		},
		Reminders: checkly.Reminders{
			Interval: 5,
		},
		SSLCertificates: checkly.SSLCertificates{
			Enabled:        false,
			AlertThreshold: 3,
		},
	}

	check := checkly.Check{
		Name:                 apiCheck.Name,
		Type:                 checkly.TypeAPI,
		Frequency:            5,
		DegradedResponseTime: 5000,
		MaxResponseTime:      15000,
		Activated:            false, // deactivated for development
		Muted:                true,  // muted for development
		ShouldFail:           false,
		DoubleCheck:          false,
		SSLCheck:             false,
		LocalSetupScript:     "",
		LocalTearDownScript:  "",
		Locations: []string{
			"eu-west-1",
			"ap-northeast-2",
		},
		Tags: []string{
			"foo",
			"bar",
		},
		AlertSettings:          alertSettings,
		UseGlobalAlertSettings: false,
		Request: checkly.Request{
			Method:  http.MethodGet,
			URL:     apiCheck.Spec.Endpoint,
			Headers: []checkly.KeyValue{
				// {
				// 	Key:   "X-Test",
				// 	Value: "foo",
				// },
			},
			QueryParameters: []checkly.KeyValue{
				// {
				// 	Key:   "query",
				// 	Value: "foo",
				// },
			},
			Assertions: []checkly.Assertion{
				{
					Source:     checkly.StatusCode,
					Comparison: checkly.Equals,
					Target:     apiCheck.Spec.Success,
				},
			},
			Body:     "",
			BodyType: "NONE",
		},
	}

	client, ctx, cancel := checklyClient()
	defer cancel()

	checklyGet, err := client.Update(ctx, apiCheck.Status.ID, check)
	if err != nil {
		return
	}

	log.Log.Info("Updated check", "check ID", checklyGet.ID)

	return
}

func delete(ID string) (err error) {

	client, ctx, cancel := checklyClient()
	defer cancel()

	err = client.Delete(ctx, ID)

	return
}
