package controllers

import (
	"context"

	"go.uber.org/zap"

	"github.com/go-logr/zapr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/slack-go/slack"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	// +kubebuilder:scaffold:imports
)

// MockSlackConnector an implementation of the SlackConnector
type MockSlackConnector struct {
}

// NewDefaultSlackConnector generates a Slack connector based on slack-go
func NewMockSlackConnector() SlackConnector {
	return &MockSlackConnector{}
}

func (mc *MockSlackConnector) GetAPIwithChannel(cm *ChatManager, chat *mysqlv1alpha1.Chat) (*slack.Client, string, error) {
	if chat != nil && chat.Spec.Slack.Channel == "doesnotexist" {
		return nil, "doesnotexist", ErrChannelNotFound
	}
	return nil, "mysql", nil
}

func (mc *MockSlackConnector) PostMessage(api *slack.Client, channel string, message string) error {
	return nil
}

var _ = Describe("Chat Controller", func() {
	It("Create a chat resource with success", func() {
		ctx := context.Background()
		chat := mysqlv1alpha1.Chat{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "chat-succeed-",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.ChatSpec{
				Slack: mysqlv1alpha1.SlackSpec{
					Channel: "mysql",
					Token:   "xoxb-...",
				},
			},
		}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &ChatReconciler{
			Client:    k8sClient,
			Log:       zapr.NewLogger(zapLog),
			Scheme:    scheme.Scheme,
			Connector: NewMockSlackConnector(),
			Chats:     map[string]*slack.Client{},
		}
		Expect(k8sClient.Create(ctx, &chat)).To(Succeed())

		chatName := types.NamespacedName{Namespace: chat.Namespace, Name: chat.Name}
		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: chatName})).To(Equal(ctrl.Result{Requeue: false}))
		response := mysqlv1alpha1.Chat{}
		Expect(k8sClient.Get(ctx, chatName, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.ChatPending), "Expected reconcile to change the status to Pending")

		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: chatName})).To(Equal(ctrl.Result{Requeue: false}))
		Expect(k8sClient.Get(ctx, chatName, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.ChatSucceeded), "Expected reconcile to change the status to Succeed")
	})

	It("Create a chat resource with failure", func() {
		ctx := context.Background()
		chat := mysqlv1alpha1.Chat{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "chat-failed-",
				Namespace:    "default",
			},
			Spec: mysqlv1alpha1.ChatSpec{
				Slack: mysqlv1alpha1.SlackSpec{
					Channel: "doesnotexist",
					Token:   "xoxb-...",
				},
			},
		}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &ChatReconciler{
			Client:    k8sClient,
			Log:       zapr.NewLogger(zapLog),
			Scheme:    scheme.Scheme,
			Connector: NewMockSlackConnector(),
			Chats:     map[string]*slack.Client{},
		}
		Expect(k8sClient.Create(ctx, &chat)).To(Succeed())

		chatName := types.NamespacedName{Namespace: chat.Namespace, Name: chat.Name}
		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: chatName})).To(Equal(ctrl.Result{Requeue: false}))
		response := mysqlv1alpha1.Chat{}
		Expect(k8sClient.Get(ctx, chatName, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.ChatPending), "Expected reconcile to change the status to Pending")

		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: chatName})).To(Equal(ctrl.Result{Requeue: false}))
		Expect(k8sClient.Get(ctx, chatName, &response)).To(Succeed())
		Expect(response.Status.Reason).To(Equal(mysqlv1alpha1.ChatSlackChannelError), "Expected reconcile to change the status to ChannelNotFound")
	})
})
