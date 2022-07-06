package checkly

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/checkly/checkly-go-sdk"
	checklyv1alpha1 "github.com/checkly/checkly-operator/apis/checkly/v1alpha1"
	external "github.com/checkly/checkly-operator/external/checkly"
)

// ApiCheckReconciler reconciles a ApiCheck object
type ApiCheckReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	ApiClient checkly.Client
}

//+kubebuilder:rbac:groups=k8s.checklyhq.com,resources=apichecks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.checklyhq.com,resources=apichecks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.checklyhq.com,resources=apichecks/finalizers,verbs=update
//+kubebuilder:rbac:groups=k8s.checklyhq.com,resources=groups,verbs=get;list

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
	logger := log.FromContext(ctx)

	apiCheckFinalizer := "k8s.checklyhq.com/finalizer"

	apiCheck := &checklyv1alpha1.ApiCheck{}

	// ////////////////////////////////
	// Delete Logic
	// ///////////////////////////////
	err := r.Get(ctx, req.NamespacedName, apiCheck)
	if err != nil {
		if errors.IsNotFound(err) {
			// The resource has been deleted
			logger.Info("Deleted", "checkly ID", apiCheck.Status.ID, "endpoint", apiCheck.Spec.Endpoint, "name", apiCheck.Name)
			return ctrl.Result{}, nil
		}
		// Error reading the object
		logger.Error(err, "can't read the object")
		return ctrl.Result{}, nil
	}

	if apiCheck.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(apiCheck, apiCheckFinalizer) {
			logger.Info("Finalizer is present, trying to delete Checkly check", "checkly ID", apiCheck.Status.ID)
			err := external.Delete(apiCheck.Status.ID, r.ApiClient)
			if err != nil {
				logger.Error(err, "Failed to delete checkly API check")
				return ctrl.Result{}, err
			}

			logger.Info("Successfully deleted checkly API check", "checkly ID", apiCheck.Status.ID)

			controllerutil.RemoveFinalizer(apiCheck, apiCheckFinalizer)
			err = r.Update(ctx, apiCheck)
			if err != nil {
				return ctrl.Result{}, err
			}
			logger.Info("Successfully deleted finalizer")
		}
		return ctrl.Result{}, nil
	}

	// Object found, let's do something with it. It's either updated, or it's new.
	logger.Info("Object found", "endpoint", apiCheck.Spec.Endpoint)

	// /////////////////////////////
	// Finalizer logic
	// ////////////////////////////
	if !controllerutil.ContainsFinalizer(apiCheck, apiCheckFinalizer) {
		controllerutil.AddFinalizer(apiCheck, apiCheckFinalizer)
		err = r.Update(ctx, apiCheck)
		if err != nil {
			logger.Error(err, "Failed to update ApiCheck status")
			return ctrl.Result{}, err
		}
		logger.Info("Added finalizer", "checkly ID", apiCheck.Status.ID, "endpoint", apiCheck.Spec.Endpoint)
		return ctrl.Result{}, nil
	}

	// /////////////////////////////
	// Lookup group ID
	// ////////////////////////////
	group := &checklyv1alpha1.Group{}
	err = r.Get(ctx, types.NamespacedName{Name: apiCheck.Spec.Group, Namespace: apiCheck.Namespace}, group)
	if err != nil {
		if errors.IsNotFound(err) {
			// The resource has been deleted
			logger.Info("Group not found, probably deleted or does not exist", "name", apiCheck.Spec.Group)
			return ctrl.Result{}, err
		}
		// Error reading the object
		logger.Error(err, "can't read the group object")
		return ctrl.Result{}, err
	}

	if group.Status.ID == 0 {
		logger.Info("Group ID has not been populated, we're too quick, requeining for retry", "group name", apiCheck.Spec.Group)
		return ctrl.Result{Requeue: true}, nil
	}

	// Create internal Check type
	internalCheck := external.Check{
		Name:            apiCheck.Name,
		Namespace:       apiCheck.Namespace,
		Frequency:       apiCheck.Spec.Frequency,
		MaxResponseTime: apiCheck.Spec.Frequency,
		Locations:       apiCheck.Spec.Locations,
		Endpoint:        apiCheck.Spec.Endpoint,
		SuccessCode:     apiCheck.Spec.Success,
		ID:              apiCheck.Status.ID,
		GroupID:         group.Status.ID,
		Muted:           apiCheck.Spec.Muted,
	}

	// /////////////////////////////
	// Update logic
	// ////////////////////////////

	// Determine if it's a new object or if it's an update to an existing object
	if apiCheck.Status.ID != "" {
		// Existing object, we need to update it
		logger.Info("Existing object, with ID", "checkly ID", apiCheck.Status.ID, "endpoint", apiCheck.Spec.Endpoint)
		err := external.Update(internalCheck, r.ApiClient)
		// err :=
		if err != nil {
			logger.Error(err, "Failed to update the checkly check")
			return ctrl.Result{}, err
		}
		logger.Info("Updated checkly check", "checkly ID", apiCheck.Status.ID)
		return ctrl.Result{}, nil
	}

	// /////////////////////////////
	// Create logic
	// ////////////////////////////

	checklyID, err := external.Create(internalCheck, r.ApiClient)
	if err != nil {
		logger.Error(err, "Failed to create checkly alert")
		return ctrl.Result{}, nil
	}

	// Update the custom resource Status with the returned ID

	apiCheck.Status.ID = checklyID
	apiCheck.Status.GroupID = group.Status.ID
	err = r.Status().Update(ctx, apiCheck)
	if err != nil {
		logger.Error(err, "Failed to update ApiCheck status")
		return ctrl.Result{}, err
	}
	logger.Info("New checkly check created with", "checkly ID", apiCheck.Status.ID, "endpoint", apiCheck.Spec.Endpoint)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApiCheckReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&checklyv1alpha1.ApiCheck{}).
		Complete(r)
}
