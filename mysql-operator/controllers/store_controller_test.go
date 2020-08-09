/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"go.uber.org/zap"

	"github.com/blaqkube/mysql-operator/mysql-operator/helpers"
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

var _ = Describe("Store Controller", func() {
	It("New Store and S3 Okay", func() {
		ctx := context.Background()
		store := mysqlv1alpha1.Store{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-store",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.StoreSpec{
				S3Backup: &mysqlv1alpha1.S3BackupInfo{
					Bucket: "pong",
				},
			},
		}

		Expect(k8sClient.Create(ctx, &store)).To(Succeed())

		name := types.NamespacedName{Namespace: store.Namespace, Name: store.Name}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &StoreReconciler{
			Client:      k8sClient,
			Log:         zapr.NewLogger(zapLog),
			Scheme:      scheme.Scheme,
			BackupStore: helpers.NewStoreMockInitialize(),
		}
		Expect(reconcile.Reconcile(ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{Requeue: true}))

		response := mysqlv1alpha1.Store{}
		Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
		Expect(response.Status.Status).To(Equal("Pending"), "Expected reconcile to change the status to Pending")

		Expect(reconcile.Reconcile(ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{}))
		Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
		Expect(response.Status.Status).To(Equal("Success"), "Expected reconcile to change the status to Success")

	})
})
