package controllers

import (
	"context"
	"errors"

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
	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/stretchr/testify/mock"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	// +kubebuilder:scaffold:imports
)

const (
	storeMockStatusSucceed  = "succeed"
	storeMockStatusWithKeys = "keys"
	storeMockStatusFailS3   = "fail/s3"
)

// NewStorage takes a S3 connection and creates a default storage
func NewStorage(status string) *Storage {
	return &Storage{
		Status: status,
	}
}

// Storage is the default storage for S3
type Storage struct {
	mock.Mock
	Status string
}

// Push pushes a file
func (s *Storage) Push(backup *openapi.BackupRequest, filename string) error {
	switch s.Status {
	case storeMockStatusFailS3:
		return errors.New("WriteFailure")
	case storeMockStatusWithKeys:
		count := 0
		for _, v := range backup.Envs {
			if (v.Name == "AWS_ACCESS_KEY_ID" && v.Value == "AKIA") ||
				(v.Name == "AWS_SECRET_ACCESS_KEY" && v.Value == "secret...") ||
				(v.Name == "AWS_REGION" && v.Value == "us-east-1") {
				count++
			}
		}
		if count == 3 {
			return nil
		}
		return errors.New("WriteFailure")
	}
	return nil
}

// Pull pull a file from S3, using a different location if necessary
func (s *Storage) Pull(backup *openapi.BackupRequest, filename string) error {
	return nil
}

// Delete deletes a file from S3
func (s *Storage) Delete(backup *openapi.BackupRequest) error {
	return nil
}

var _ = Describe("Store Controller", func() {
	It("Create a new store and check success/failure", func() {
		ctx := context.Background()

		tests := []map[string]string{
			{
				"status": storeMockStatusSucceed,
				"result": mysqlv1alpha1.StoreCheckSucceeded,
			},
			{
				"status": storeMockStatusFailS3,
				"result": mysqlv1alpha1.StoreCheckFailed,
			},
		}

		for _, test := range tests {
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
			reconcile := &StoreReconciler{
				Client: k8sClient,
				Log:    zapr.NewLogger(zapLog),
				Scheme: scheme.Scheme,
				Storages: map[string]backend.Storage{
					"s3":        NewStorage(test["status"]),
					"blackhole": NewStorage(test["status"]),
					"gcp":       NewStorage(test["status"]),
				},
			}

			Expect(k8sClient.Create(ctx, &store)).To(Succeed())

			name := types.NamespacedName{Namespace: store.Namespace, Name: store.Name}
			Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: false}))
			response := mysqlv1alpha1.Store{}
			Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
			Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.StoreCheckRequested), "Expected reconcile to change the status to Check")

			Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: false}))
			Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
			Expect(response.Status.Reason).To(Equal(test["result"]), "Expected reconcile to change the status to the result")

			Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: false}))
		}
	})

	It("Create a store with SecretKeyRefs/succeed", func() {
		ctx := context.Background()

		secret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "store-secret1",
				Namespace: "default",
			},
			Type: corev1.SecretTypeOpaque,
			StringData: map[string]string{
				"AWS_ACCESS_KEY_ID":     "AKIA",
				"AWS_SECRET_ACCESS_KEY": "secret...",
				"AWS_REGION":            "us-east-1",
			},
		}

		store := mysqlv1alpha1.Store{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "store-",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.StoreSpec{
				Bucket: "pong",
				Envs: []corev1.EnvVar{
					{
						Name: "AWS_ACCESS_KEY_ID",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{Name: "store-secret1"},
								Key:                  "AWS_ACCESS_KEY_ID",
							},
						},
					},
					{
						Name: "AWS_SECRET_ACCESS_KEY",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{Name: "store-secret1"},
								Key:                  "AWS_SECRET_ACCESS_KEY",
							},
						},
					},
					{
						Name: "AWS_REGION",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{Name: "store-secret1"},
								Key:                  "AWS_REGION",
							},
						},
					},
				},
			},
		}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &StoreReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
			Storages: map[string]backend.Storage{
				"s3":        NewStorage(storeMockStatusWithKeys),
				"blackhole": NewStorage(storeMockStatusWithKeys),
				"gcp":       NewStorage(storeMockStatusWithKeys),
			},
		}

		Expect(k8sClient.Create(ctx, &secret)).To(Succeed())
		Expect(k8sClient.Create(ctx, &store)).To(Succeed())

		name := types.NamespacedName{Namespace: store.Namespace, Name: store.Name}
		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: false}))
		response := mysqlv1alpha1.Store{}
		Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.StoreCheckRequested), "Expected reconcile to change the status to Check")

		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: false}))
		Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.StoreCheckSucceeded), "Expected reconcile to change the status to the result")

		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: false}))
	})

	It("Create a store with SecretKeyRefs/Missing keys", func() {
		ctx := context.Background()

		secret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "store-secret2",
				Namespace: "default",
			},
			Type: corev1.SecretTypeOpaque,
			StringData: map[string]string{
				"AWS_ACCESS_KEY_ID": "AKIA",
			},
		}

		store := mysqlv1alpha1.Store{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "store-",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.StoreSpec{
				Bucket: "pong",
				Envs: []corev1.EnvVar{
					{
						Name: "AWS_ACCESS_KEY_ID",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{Name: "store-secret2"},
								Key:                  "AWS_ACCESS_KEY_ID",
							},
						},
					},
					{
						Name: "AWS_SECRET_ACCESS_KEY",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{Name: "store-secret2"},
								Key:                  "AWS_SECRET_ACCESS_KEY",
							},
						},
					},
					{
						Name: "AWS_REGION",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{Name: "store-secret2"},
								Key:                  "AWS_REGION",
							},
						},
					},
				},
			},
		}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &StoreReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
			Storages: map[string]backend.Storage{
				"s3":        NewStorage(storeMockStatusWithKeys),
				"blackhole": NewStorage(storeMockStatusWithKeys),
				"gcp":       NewStorage(storeMockStatusWithKeys),
			},
		}

		Expect(k8sClient.Create(ctx, &secret)).To(Succeed())
		Expect(k8sClient.Create(ctx, &store)).To(Succeed())

		name := types.NamespacedName{Namespace: store.Namespace, Name: store.Name}
		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: false}))
		response := mysqlv1alpha1.Store{}
		Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.StoreCheckRequested), "Expected reconcile to change the status to Check")

		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: false}))
		Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.StoreCheckFailed), "Expected reconcile to change the status to the result")

		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: false}))
	})

	It("Create a store with SecretKeyRefs/succeed", func() {
		ctx := context.Background()

		secret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "store-secret3",
				Namespace: "default",
			},
			Type: corev1.SecretTypeOpaque,
			StringData: map[string]string{
				"AWS_ACCESS_KEY_ID":     "AKIA",
				"AWS_SECRET_ACCESS_KEY": "secret...",
				"AWS_REGION":            "us-east-1",
			},
		}

		store := mysqlv1alpha1.Store{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "store-",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.StoreSpec{
				Bucket: "pong",
				Envs: []corev1.EnvVar{
					{
						Name: "AWS_ACCESS_KEY_ID",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{Name: "store-secret3"},
								Key:                  "AWS_ACCESS_KEY_ID",
							},
						},
					},
				},
			},
		}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &StoreReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
			Storages: map[string]backend.Storage{
				"s3":        NewStorage(storeMockStatusWithKeys),
				"blackhole": NewStorage(storeMockStatusWithKeys),
				"gcp":       NewStorage(storeMockStatusWithKeys),
			},
		}

		Expect(k8sClient.Create(ctx, &secret)).To(Succeed())
		Expect(k8sClient.Create(ctx, &store)).To(Succeed())

		name := types.NamespacedName{Namespace: store.Namespace, Name: store.Name}
		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: false}))
		response := mysqlv1alpha1.Store{}
		Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.StoreCheckRequested), "Expected reconcile to change the status to Check")

		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: false}))
		Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.StoreCheckFailed), "Expected reconcile to change the status to the result")

		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: false}))
	})

})
