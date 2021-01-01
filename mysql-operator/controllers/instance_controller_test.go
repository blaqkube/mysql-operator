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

	"github.com/blaqkube/mysql-operator/agent/backend"
	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	// +kubebuilder:scaffold:imports
)

var _ = Describe("Instance Controller", func() {
	It("Create instance without any restore", func() {
		ctx := context.Background()
		instance := mysqlv1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-mysql1",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.InstanceSpec{
				Database: "me",
			},
		}

		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())

		instanceName := types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &InstanceReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
			Properties: &StatefulSetProperties{
				AgentVersion: "latest",
				MySQLVersion: "8.0.22",
			},
		}
		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: instanceName})).To(Equal(ctrl.Result{}))

		response := mysqlv1alpha1.Instance{}
		Expect(k8sClient.Get(ctx, instanceName, &response)).To(Succeed())
		Expect(mysqlv1alpha1.InstanceExporterSecretCreated).
			To(Equal(response.Status.Reason), "Expected reconcile to change the status to ExporterSecretCreated")

		secret := corev1.Secret{}
		secretName := types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name + "-exporter"}
		Expect(k8sClient.Get(ctx, secretName, &secret)).To(Succeed())
		Expect(string(secret.Data[".my.cnf"])).
			Should(MatchRegexp(`\[client\].*`), "Data should match the password")

		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: instanceName})).To(Equal(ctrl.Result{}))

		response = mysqlv1alpha1.Instance{}
		Expect(k8sClient.Get(ctx, instanceName, &response)).To(Succeed())
		Expect(mysqlv1alpha1.InstanceStatefulSetCreated).
			To(Equal(response.Status.Reason), "Expected reconcile to change the status to StatefulSetCreated")
	})

	It("Create an Instance with a Store", func() {

		ctx := context.TODO()

		store := mysqlv1alpha1.Store{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "store-",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.StoreSpec{
				Bucket: "pong",
			},
		}

		zapLog, _ := zap.NewDevelopment()
		storeReconcile := &StoreReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
			Storages: map[string]backend.Storage{
				"s3":        NewStorage(storeMockStatusSucceed),
				"blackhole": NewStorage(storeMockStatusSucceed),
				"gcp":       NewStorage(storeMockStatusSucceed),
			},
		}

		Expect(k8sClient.Create(ctx, &store)).To(Succeed())

		storeName := types.NamespacedName{Namespace: store.Namespace, Name: store.Name}
		Expect(storeReconcile.Reconcile(ctx, ctrl.Request{NamespacedName: storeName})).To(Equal(ctrl.Result{Requeue: false}))
		storeResponse := mysqlv1alpha1.Store{}
		Expect(k8sClient.Get(ctx, storeName, &storeResponse)).To(Succeed())
		Expect(storeResponse.Status.Reason).To(Equal(mysqlv1alpha1.StoreCheckRequested), "Expected reconcile to change the status to Check")

		Expect(storeReconcile.Reconcile(ctx, ctrl.Request{NamespacedName: storeName})).To(Equal(ctrl.Result{Requeue: false}))
		Expect(k8sClient.Get(ctx, storeName, &storeResponse)).To(Succeed())
		Expect(storeResponse.Status.Reason).To(Equal(mysqlv1alpha1.StoreCheckSucceeded), "Expected reconcile to change the status to the result")

		Expect(storeReconcile.Reconcile(ctx, ctrl.Request{NamespacedName: storeName})).To(Equal(ctrl.Result{Requeue: false}))

		instance := mysqlv1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "instance-",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.InstanceSpec{
				Database: "me",
				Restore: mysqlv1alpha1.RestoreSpec{
					Store:    store.Name,
					Location: "/location/backup01.sql",
				},
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())

		instanceName := types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name}

		instanceReconcile := &InstanceReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
			Properties: &StatefulSetProperties{
				AgentVersion: "latest",
				MySQLVersion: "8.0.22",
			},
		}
		Expect(instanceReconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: instanceName})).To(Equal(ctrl.Result{}))

		instanceResponse := mysqlv1alpha1.Instance{}
		Expect(k8sClient.Get(ctx, instanceName, &instanceResponse)).To(Succeed())
		Expect(mysqlv1alpha1.InstanceExporterSecretCreated).
			To(Equal(instanceResponse.Status.Reason), "Expected reconcile to change the status to ExporterSecretCreated")

		secret := corev1.Secret{}
		secretName := types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name + "-exporter"}
		Expect(k8sClient.Get(ctx, secretName, &secret)).To(Succeed())
		Expect(string(secret.Data[".my.cnf"])).
			Should(MatchRegexp(`\[client\].*`), "Data should match the password")

		Expect(instanceReconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: instanceName})).To(Equal(ctrl.Result{}))

		instanceResponse = mysqlv1alpha1.Instance{}
		Expect(k8sClient.Get(ctx, instanceName, &instanceResponse)).To(Succeed())
		Expect(mysqlv1alpha1.InstanceStatefulSetCreated).
			To(Equal(instanceResponse.Status.Reason), "Expected reconcile to change the status to StatefulSetCreated")
	})

	It("Create an Instance with later Store", func() {

		storeName := types.NamespacedName{Namespace: "default", Name: "store-later-1"}
		ctx := context.TODO()
		zapLog, _ := zap.NewDevelopment()

		instance := mysqlv1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "instance-",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.InstanceSpec{
				Database: "me",
				Restore: mysqlv1alpha1.RestoreSpec{
					Store:    storeName.Name,
					Location: "/location/backup01.sql",
				},
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())

		instanceName := types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name}

		instanceReconcile := &InstanceReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
			Properties: &StatefulSetProperties{
				AgentVersion: "latest",
				MySQLVersion: "8.0.22",
			},
		}
		Expect(instanceReconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: instanceName})).To(Equal(ctrl.Result{}))

		instanceResponse := mysqlv1alpha1.Instance{}
		Expect(k8sClient.Get(ctx, instanceName, &instanceResponse)).To(Succeed())
		Expect(mysqlv1alpha1.InstanceExporterSecretCreated).
			To(Equal(instanceResponse.Status.Reason), "Expected reconcile to change the status to ExporterSecretCreated")

		secret := corev1.Secret{}
		secretName := types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name + "-exporter"}
		Expect(k8sClient.Get(ctx, secretName, &secret)).To(Succeed())
		Expect(string(secret.Data[".my.cnf"])).
			Should(MatchRegexp(`\[client\].*`), "Data should match the password")

		Expect(instanceReconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: instanceName})).To(Equal(ctrl.Result{}))

		instanceResponse = mysqlv1alpha1.Instance{}
		Expect(k8sClient.Get(ctx, instanceName, &instanceResponse)).To(Succeed())
		Expect(mysqlv1alpha1.InstanceStoreInaccessible).
			To(Equal(instanceResponse.Status.Reason), "Expected reconcile to change the status to StatefulSetCreated")

		store := mysqlv1alpha1.Store{
			ObjectMeta: metav1.ObjectMeta{
				Name:      storeName.Name,
				Namespace: storeName.Namespace,
			},
			Spec: mysqlv1alpha1.StoreSpec{
				Bucket: "pong",
			},
		}

		storeReconcile := &StoreReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
			Storages: map[string]backend.Storage{
				"s3":        NewStorage(storeMockStatusSucceed),
				"blackhole": NewStorage(storeMockStatusSucceed),
				"gcp":       NewStorage(storeMockStatusSucceed),
			},
		}

		Expect(k8sClient.Create(ctx, &store)).To(Succeed())

		Expect(storeReconcile.Reconcile(ctx, ctrl.Request{NamespacedName: storeName})).To(Equal(ctrl.Result{Requeue: false}))
		storeResponse := mysqlv1alpha1.Store{}
		Expect(k8sClient.Get(ctx, storeName, &storeResponse)).To(Succeed())
		Expect(storeResponse.Status.Reason).To(Equal(mysqlv1alpha1.StoreCheckRequested), "Expected reconcile to change the status to Check")

		Expect(storeReconcile.Reconcile(ctx, ctrl.Request{NamespacedName: storeName})).To(Equal(ctrl.Result{Requeue: false}))
		Expect(k8sClient.Get(ctx, storeName, &storeResponse)).To(Succeed())
		Expect(storeResponse.Status.Reason).To(Equal(mysqlv1alpha1.StoreCheckSucceeded), "Expected reconcile to change the status to the result")

		Expect(storeReconcile.Reconcile(ctx, ctrl.Request{NamespacedName: storeName})).To(Equal(ctrl.Result{Requeue: false}))

		Expect(instanceReconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: instanceName})).To(Equal(ctrl.Result{}))

		instanceResponse = mysqlv1alpha1.Instance{}
		Expect(k8sClient.Get(ctx, instanceName, &instanceResponse)).To(Succeed())
		Expect(mysqlv1alpha1.InstanceStatefulSetCreated).
			To(Equal(instanceResponse.Status.Reason), "Expected reconcile to change the status to StatefulSetCreated")
	})

})
