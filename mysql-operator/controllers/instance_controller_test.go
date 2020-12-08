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
	"github.com/blaqkube/mysql-operator/mysql-operator/helpers"
	// +kubebuilder:scaffold:imports
)

var _ = Describe("Instance Controller", func() {
	It("Instance Okay", func() {
		ctx := context.Background()
		instance := mysqlv1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-mysql1",
				Namespace: "default",
			},
			Spec: mysqlv1alpha1.InstanceSpec{
				Database: "me",
			},
		}

		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())

		name := types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &InstanceReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
			Properties: StatefulSetProperties{
				AgentVersion: "latest",
				MySQLVersion: "8.0.21",
			},
		}
		Expect(reconcile.Reconcile(ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{}))

		response := mysqlv1alpha1.Instance{}
		Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
		Expect("Success").To(Equal(response.Status.LastCondition), "Expected reconcile to change the status to Success")
	})

	It("Instance with existing Store", func() {
		ctx := context.Background()
		instance := mysqlv1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-mysql2",
				Namespace: "default",
			},
			Spec: mysqlv1alpha1.InstanceSpec{
				Database: "me",
				Restore: mysqlv1alpha1.RestoreSpec{
					Store: "existing-store",
				},
			},
		}

		store := mysqlv1alpha1.Store{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "existing-store",
				Namespace: "default",
			},
			Spec: mysqlv1alpha1.StoreSpec{
				S3Backup: &mysqlv1alpha1.S3BackupInfo{
					Bucket: "pong",
				},
			},
		}

		Expect(k8sClient.Create(ctx, &store)).To(Succeed())

		storeName := types.NamespacedName{Namespace: store.Namespace, Name: store.Name}

		zapLog, _ := zap.NewDevelopment()
		storeReconcile := &StoreReconciler{
			Client:      k8sClient,
			Log:         zapr.NewLogger(zapLog),
			Scheme:      scheme.Scheme,
			BackupStore: helpers.NewStoreMockInitialize(),
		}
		Expect(storeReconcile.Reconcile(ctrl.Request{NamespacedName: storeName})).To(Equal(ctrl.Result{Requeue: true}))

		storeResponse := mysqlv1alpha1.Store{}
		Expect(k8sClient.Get(ctx, storeName, &storeResponse)).To(Succeed())
		Expect(storeResponse.Status.LastCondition).To(Equal("Pending"), "Expected reconcile to change the status to Pending")

		Expect(storeReconcile.Reconcile(ctrl.Request{NamespacedName: storeName})).To(Equal(ctrl.Result{}))
		Expect(k8sClient.Get(ctx, storeName, &storeResponse)).To(Succeed())
		Expect(storeResponse.Status.LastCondition).To(Equal("Success"), "Expected reconcile to change the status to Success")

		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())

		name := types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name}

		instanceReconcile := &InstanceReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
			Properties: StatefulSetProperties{
				AgentVersion: "",
				MySQLVersion: "",
			},
		}
		Expect(instanceReconcile.Reconcile(ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{}))

		instanceResponse := mysqlv1alpha1.Instance{}
		Expect(k8sClient.Get(ctx, name, &instanceResponse)).To(Succeed())
		Expect("Success").To(Equal(instanceResponse.Status.LastCondition), "Expected reconcile to change the status to Scheduling")
	})

	It("Instance with Unexisting Store", func() {
		ctx := context.Background()
		instance := mysqlv1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-mysql3",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.InstanceSpec{
				Database: "me",
				Restore: mysqlv1alpha1.RestoreSpec{
					Store: "missing-store",
				},
			},
		}

		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())

		name := types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &InstanceReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
			Properties: StatefulSetProperties{
				AgentVersion: "",
				MySQLVersion: "",
			},
		}
		Expect(reconcile.Reconcile(ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: true, RequeueAfter: time.Duration(30 * time.Second)}))

		response := mysqlv1alpha1.Instance{}
		Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
		Expect("Waiting for store").To(Equal(response.Status.LastCondition), "Expected reconcile to change the status to Scheduling")
	})
})
