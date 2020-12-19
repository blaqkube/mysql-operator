package controllers

import (
	"context"

	"go.uber.org/zap"

	"github.com/go-logr/zapr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	// +kubebuilder:scaffold:imports
)

var _ = Describe("Store Controller", func() {
	It("New Store and S3 Okay", func() {
		ctx := context.Background()
		store := mysqlv1alpha1.Store{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-store",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.StoreSpec{
				Bucket: "pong",
			},
		}

		Expect(k8sClient.Create(ctx, &store)).To(Succeed())

		name := types.NamespacedName{Namespace: store.Namespace, Name: store.Name}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &StoreReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
		}
		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: true}))

		response := mysqlv1alpha1.Store{}
		Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal("Check"), "Expected reconcile to change the status to Check")

		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{}))
		Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal("Success"), "Expected reconcile to change the status to Success")

	})

	It("New Store fails", func() {
		ctx := context.Background()
		store := mysqlv1alpha1.Store{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-store",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.StoreSpec{
				Bucket: "fail",
				Env:    []corev1.EnvVar{},
			},
		}

		Expect(k8sClient.Create(ctx, &store)).To(Succeed())

		name := types.NamespacedName{Namespace: store.Namespace, Name: store.Name}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &StoreReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
		}
		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: true}))

		response := mysqlv1alpha1.Store{}
		Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal("Check"), "Expected reconcile to change the status to Check")

		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{}))
		Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal("Success"), "Expected reconcile to change the status to AccessDenied")
	})

	It("Write Store fails", func() {
		ctx := context.Background()
		store := mysqlv1alpha1.Store{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-store",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.StoreSpec{
				Bucket: "fail",
				Env:    []corev1.EnvVar{},
			},
		}

		Expect(k8sClient.Create(ctx, &store)).To(Succeed())

		name := types.NamespacedName{Namespace: store.Namespace, Name: store.Name}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &StoreReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
		}
		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: true}))

		response := mysqlv1alpha1.Store{}
		Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal("Check"), "Expected reconcile to change the status to Check")

		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{}))
		Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal("Success"), "Expected reconcile to change the status to AccessDenied")
	})
})
