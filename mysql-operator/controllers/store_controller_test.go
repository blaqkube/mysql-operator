package controllers

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/go-logr/zapr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/stretchr/testify/mock"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	// +kubebuilder:scaffold:imports
)

const (
	statusSucceed = "succeed"
	statusFailS3  = "fail/s3"
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
	case statusFailS3:
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
				"status": statusSucceed,
				"result": mysqlv1alpha1.StateCheckSucceeded,
			},
			{
				"status": statusFailS3,
				"result": mysqlv1alpha1.StateCheckFailed,
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
				Client:  k8sClient,
				Log:     zapr.NewLogger(zapLog),
				Scheme:  scheme.Scheme,
				Storage: NewStorage(test["status"]),
			}

			Expect(k8sClient.Create(ctx, &store)).To(Succeed())

			name := types.NamespacedName{Namespace: store.Namespace, Name: store.Name}
			Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: false}))
			response := mysqlv1alpha1.Store{}
			Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
			Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.StateCheckRequested), "Expected reconcile to change the status to Check")

			Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: false}))
			Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
			Expect(response.Status.Reason).To(Equal(test["result"]), "Expected reconcile to change the status to the result")

			Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: false}))
		}
	})

})
