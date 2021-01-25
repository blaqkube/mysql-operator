package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/slack-go/slack"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
)

// ChatReconciler reconciles a Chat object
type ChatReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Chats  map[string]*slack.Client
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=chats,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=chats/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=chats/finalizers,verbs=update

// Reconcile implement the reconciliation loop for chats
func (r *ChatReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("chat", req.NamespacedName)
	log.Info("Running a reconcile loop")

	chat := mysqlv1alpha1.Chat{}
	if err := r.Get(ctx, req.NamespacedName, &chat); err != nil {
		log.Info("Unable to fetch chat from kubernetes")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	cm := &ChatManager{
		Context:     ctx,
		Reconciler:  r,
		TimeManager: NewTimeManager(),
	}

	if chat.Status.Reason == "" {
		chat.Status.Reason = mysqlv1alpha1.ChatPending
		chat.Status.Ready = metav1.ConditionUnknown

		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionUnknown,
			LastTransitionTime: metav1.Now(),
			Reason:             mysqlv1alpha1.ChatPending,
			Message:            "A new check has been requested",
		}
		return cm.setChatCondition(&chat, condition)
	}

	// TODO: Reconciler should be able to
	// - detect a change in the ConfigMap or Secret and reload the associated data
	// - Retry on regular basis in the event of a failure
	// - Update chat status when stchat chatore moves to success
	if chat.Status.Reason == mysqlv1alpha1.ChatPending {
		api, channel, err := cm.GetAPIwithChannel(&chat)
		condition := metav1.Condition{}
		if err != nil {
			switch err {
			case ErrChannelNotFound:
				chat.Status.Reason = mysqlv1alpha1.ChatSlackChannelError
				chat.Status.Ready = metav1.ConditionFalse

				condition = metav1.Condition{
					Type:               "available",
					Status:             metav1.ConditionUnknown,
					LastTransitionTime: metav1.Now(),
					Reason:             mysqlv1alpha1.ChatSlackChannelError,
					Message:            "Could not find Slack Channel",
				}
			case ErrChatConnectionFailed:
				chat.Status.Reason = mysqlv1alpha1.ChatSlackConnectionError
				chat.Status.Ready = metav1.ConditionFalse

				condition = metav1.Condition{
					Type:               "available",
					Status:             metav1.ConditionUnknown,
					LastTransitionTime: metav1.Now(),
					Reason:             mysqlv1alpha1.ChatSlackConnectionError,
					Message:            "Could not connect to Slack",
				}
			}
			return cm.setChatCondition(&chat, condition)
		}
		api.PostMessage(channel, slack.MsgOptionText(
			fmt.Sprintf("Blaqkube Chat %s/%s succeeded", chat.Namespace, chat.Name),
			false,
		))
		chat.Status.Reason = mysqlv1alpha1.ChatSucceeded
		chat.Status.Ready = metav1.ConditionTrue
		condition = metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
			Reason:             mysqlv1alpha1.ChatSucceeded,
			Message:            fmt.Sprintf("A message has been sent to %s", chat.Spec.Slack.Channel),
		}
		return cm.setChatCondition(&chat, condition)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChatReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.Chat{}).
		Complete(r)
}
