package controllers

import (
	"context"

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

var _ = Describe("Backup Controller", func() {
	It("Create a backup without any store/database", func() {
		ctx := context.Background()
		backup := mysqlv1alpha1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "backup-",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.BackupSpec{
				Store:    "store",
				Instance: "instance",
			},
		}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &BackupReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
		}
		Expect(k8sClient.Create(ctx, &backup)).To(Succeed())

		backupName := types.NamespacedName{Namespace: backup.Namespace, Name: backup.Name}
		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: backupName})).To(Equal(ctrl.Result{Requeue: false}))
		response := mysqlv1alpha1.Backup{}
		Expect(k8sClient.Get(ctx, backupName, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.BackupStoreAccessError), "Expected reconcile to change the status to Check")

	})

})
