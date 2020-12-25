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
)

var _ = Describe("Database Controller", func() {
	It("Create a database without any instance", func() {
		ctx := context.Background()
		database := mysqlv1alpha1.Database{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "store-",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.DatabaseSpec{
				Name:     "ping",
				Instance: "pong",
			},
		}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &DatabaseReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
		}

		Expect(k8sClient.Create(ctx, &database)).To(Succeed())

		databaseName := types.NamespacedName{Namespace: database.Namespace, Name: database.Name}
		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: databaseName})).To(Equal(ctrl.Result{Requeue: false}))
		response := mysqlv1alpha1.Database{}
		Expect(k8sClient.Get(ctx, databaseName, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.DatabaseInstanceAccessError), "Expected reconcile to change the status to Check")

	})

})
