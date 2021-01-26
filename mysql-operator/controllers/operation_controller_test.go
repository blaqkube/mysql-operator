package controllers

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/go-logr/zapr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	// +kubebuilder:scaffold:imports
)

var _ = Describe("Operation Controller", func() {
	It("Create a NoOp operation and run in immediate mode", func() {
		ctx := context.Background()
		operation := mysqlv1alpha1.Operation{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "noop-",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.OperationSpec{
				Instance: "blue",
				Type:     mysqlv1alpha1.OperationTypeNoop,
				Mode:     mysqlv1alpha1.OperationModeImmediate,
			},
		}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &OperationReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
		}
		Expect(k8sClient.Create(ctx, &operation)).To(Succeed())

		operationName := types.NamespacedName{Namespace: operation.Namespace, Name: operation.Name}
		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: operationName})).To(Equal(ctrl.Result{Requeue: false}))
		response := mysqlv1alpha1.Operation{}
		Expect(k8sClient.Get(ctx, operationName, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.OperationRequested), "Expected reconcile to change the status to Pending")

		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: operationName})).To(Equal(ctrl.Result{Requeue: false}))
		response = mysqlv1alpha1.Operation{}
		Expect(k8sClient.Get(ctx, operationName, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.OperationSucceeded), "Expected reconcile to change the status to Succeeded")
	})

	It("Create a NoOp operation and run in maintenance mode", func() {
		ctx := context.Background()
		operation := mysqlv1alpha1.Operation{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "noop-",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.OperationSpec{
				Instance: "instance",
				Type:     mysqlv1alpha1.OperationTypeNoop,
				Mode:     mysqlv1alpha1.OperationModeMaintenance,
			},
		}

		instance := mysqlv1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "instance",
				Namespace: "default",
			},
			Spec: mysqlv1alpha1.InstanceSpec{
				Database: "blue",
			},
		}

		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())

		zapLog, _ := zap.NewDevelopment()
		reconcile := &OperationReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
		}
		Expect(k8sClient.Create(ctx, &operation)).To(Succeed())

		operationName := types.NamespacedName{Namespace: operation.Namespace, Name: operation.Name}
		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: operationName})).To(Equal(ctrl.Result{Requeue: false}))
		response := mysqlv1alpha1.Operation{}
		Expect(k8sClient.Get(ctx, operationName, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.OperationPending), "Expected reconcile to change the status to Pending")

		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: operationName})).To(Equal(ctrl.Result{Requeue: true, RequeueAfter: time.Second}))
		response = mysqlv1alpha1.Operation{}
		Expect(k8sClient.Get(ctx, operationName, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.OperationPending), "Expected reconcile to change the status to Pending")

		instance.Status.MaintenanceMode = true
		Expect(k8sClient.Status().Update(ctx, &instance)).To(Succeed())

		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: operationName})).To(Equal(ctrl.Result{Requeue: false}))
		response = mysqlv1alpha1.Operation{}
		Expect(k8sClient.Get(ctx, operationName, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.OperationRequested), "Expected reconcile to change the status to Requested")

		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: operationName})).To(Equal(ctrl.Result{Requeue: false}))
		response = mysqlv1alpha1.Operation{}
		Expect(k8sClient.Get(ctx, operationName, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.OperationSucceeded), "Expected reconcile to change the status to Succeeded")
	})
})
