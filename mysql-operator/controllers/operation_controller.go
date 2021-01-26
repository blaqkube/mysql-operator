package controllers

import (
	"context"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
)

// OperationReconciler reconciles a Operation object
type OperationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=operations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=operations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=operations/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *OperationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("operation", req.NamespacedName)
	log.Info("Running a reconcile loop")

	operation := mysqlv1alpha1.Operation{}
	if err := r.Get(ctx, req.NamespacedName, &operation); err != nil {
		log.Info("Unable to fetch operation from kubernetes")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	om := &OperationManager{
		Context:     ctx,
		Reconciler:  r,
		TimeManager: NewTimeManager(),
	}

	if operation.Status.Reason == "" {
		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             mysqlv1alpha1.OperationPending,
			Message:            "The operation is waiting for next Maintenance Window",
		}
		if operation.Spec.Mode == mysqlv1alpha1.OperationModeImmediate {
			condition.Reason = mysqlv1alpha1.OperationRequested
			condition.Message = "The operation will be started"
		}
		return om.setOperationCondition(&operation, condition)
	}

	// TODO: Reconciler should be able to
	// - detect a change in the ConfigMap or Secret and reload the associated data
	// - Retry on regular basis in the event of a failure
	// - Update chat status when stchat chatore moves to success
	if operation.Status.Reason == mysqlv1alpha1.OperationRequested {
		switch operation.Spec.Type {
		case mysqlv1alpha1.OperationTypeNoop:
			om.NoOp()
		}
		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
			Reason:             mysqlv1alpha1.OperationSucceeded,
			Message:            "The operation has been executed with success",
		}
		return om.setOperationCondition(&operation, condition)
	}

	if operation.Status.Reason == mysqlv1alpha1.OperationPending {
		instance := mysqlv1alpha1.Instance{}
		i := types.NamespacedName{
			Namespace: operation.Namespace,
			Name:      operation.Spec.Instance,
		}
		if err := r.Get(ctx, i, &instance); err != nil {
			log.Info("Unable to fetch operation from kubernetes")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		if instance.Status.MaintenanceMode == true {
			condition := metav1.Condition{
				Type:               "available",
				Status:             metav1.ConditionFalse,
				LastTransitionTime: metav1.Now(),
				Reason:             mysqlv1alpha1.OperationRequested,
				Message:            "The operation will be started",
			}
			return om.setOperationCondition(&operation, condition)
		}
		c := len(operation.Status.Conditions) - 1
		d := om.TimeManager.Next(operation.Status.Conditions[c].LastTransitionTime.Time)
		return ctrl.Result{Requeue: true, RequeueAfter: d}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OperationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.Operation{}).
		Complete(r)
}
