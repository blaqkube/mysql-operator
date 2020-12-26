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

var _ = Describe("User Controller", func() {
	It("Create a user without any instance", func() {
		ctx := context.Background()
		user := mysqlv1alpha1.User{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "store-",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.UserSpec{
				Username: "ping",
				Password: "pong",
				Instance: "pong",
			},
		}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &UserReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
		}

		Expect(k8sClient.Create(ctx, &user)).To(Succeed())

		userName := types.NamespacedName{Namespace: user.Namespace, Name: user.Name}
		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: userName})).To(Equal(ctrl.Result{Requeue: false}))
		response := mysqlv1alpha1.User{}
		Expect(k8sClient.Get(ctx, userName, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.UserInstanceAccessError), "Expected reconcile to change the status to Check")

	})

})
