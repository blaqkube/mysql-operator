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

var _ = Describe("Grant Controller", func() {
	It("Create a grant without any user/database", func() {
		ctx := context.Background()
		user := mysqlv1alpha1.Grant{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "i1-d1-u1-",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.GrantSpec{
				User:       "ping",
				Database:   "pong",
				AccessMode: "readOnly",
			},
		}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &GrantReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
		}
		Expect(k8sClient.Create(ctx, &user)).To(Succeed())

		userName := types.NamespacedName{Namespace: user.Namespace, Name: user.Name}
		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: userName})).To(Equal(ctrl.Result{Requeue: false}))
		response := mysqlv1alpha1.Grant{}
		Expect(k8sClient.Get(ctx, userName, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.GrantUserAccessError), "Expected reconcile to change the status to Check")

	})

})
