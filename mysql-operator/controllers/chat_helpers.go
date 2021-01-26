package controllers

import (
	"context"
	"errors"
	"fmt"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	"github.com/slack-go/slack"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	maxChatConditions = 10
)

var (
	// ErrChannelNotFound is reported when the channel does not exist
	ErrChannelNotFound = errors.New("ChannelNotFound")

	// ErrChatConnectionFailed is reported when the connection failed
	ErrChatConnectionFailed = errors.New("ConnectionFailed")
)

// ChatManager provides methods to manage chat subcomponents
type ChatManager struct {
	Context     context.Context
	Reconciler  *ChatReconciler
	TimeManager *TimeManager
}

func (cm *ChatManager) setChatCondition(chat *mysqlv1alpha1.Chat, condition metav1.Condition) (ctrl.Result, error) {
	if condition.Reason == chat.Status.Reason {
		c := len(chat.Status.Conditions) - 1
		d := cm.TimeManager.Next(chat.Status.Conditions[c].LastTransitionTime.Time)
		if condition.Reason != mysqlv1alpha1.ChatSucceeded {
			return ctrl.Result{Requeue: true, RequeueAfter: d}, nil
		}
		return ctrl.Result{}, nil
	}
	chat.Status.Ready = condition.Status
	chat.Status.Reason = condition.Reason
	chat.Status.Message = condition.Message
	conditions := append(chat.Status.Conditions, condition)
	if len(conditions) > maxChatConditions {
		conditions = conditions[1:]
	}
	chat.Status.Conditions = conditions
	log := cm.Reconciler.Log.WithValues("namespace", chat.Namespace, "chat", chat.Name)
	log.Info("Updating chat with new Status", "Reason", condition.Reason, "Message", condition.Message)
	if err := cm.Reconciler.Status().Update(cm.Context, chat); err != nil {
		log.Error(err, "Unable to update chat")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// getSlackToken returns the slack token for the chat
func (cm *ChatManager) getSlackToken(chat *mysqlv1alpha1.Chat) (string, error) {
	log := cm.Reconciler.Log.WithValues("chat", types.NamespacedName{Namespace: chat.Namespace, Name: chat.Name})
	if chat.Spec.Slack.Token != "" {
		return chat.Spec.Slack.Token, nil
	}

	if chat.Spec.Slack.TokenFrom != nil {
		switch {
		case chat.Spec.Slack.TokenFrom.ConfigMapKeyRef != nil:
			configMapRef := chat.Spec.Slack.TokenFrom.ConfigMapKeyRef
			namespace := chat.Namespace
			name := configMapRef.Name
			key := configMapRef.Key
			configMap := corev1.ConfigMap{}
			if err := cm.Reconciler.Get(cm.Context, types.NamespacedName{Namespace: namespace, Name: name}, &configMap); err != nil {
				log.Info(fmt.Sprintf("Read Configmap, error: %s", err), "namespace", namespace, "configmap", name)
				return "", ErrConfigMapNotFound
			}
			value, ok := configMap.Data[key]
			if !ok {
				log.Info(fmt.Sprintf("Key %s in Configmap does not exist", key), "namespace", namespace, "configmap", name)
				return "", ErrKeyNotFound
			}
			return value, nil
		case chat.Spec.Slack.TokenFrom.SecretKeyRef != nil:
			s := chat.Spec.Slack.TokenFrom.SecretKeyRef
			namespace := chat.Namespace
			name := s.Name
			key := s.Key
			secret := corev1.Secret{}
			if err := cm.Reconciler.Get(cm.Context, types.NamespacedName{Namespace: namespace, Name: name}, &secret); err != nil {
				log.Info(fmt.Sprintf("Read secret, error: %s", err), "namespace", namespace, "secret", name)
				return "", ErrSecretNotFound
			}
			valueBytes, ok := secret.Data[key]
			if !ok {
				log.Info(fmt.Sprintf("Key %s in secret does not exist", key), "namespace", namespace, "configmap", name)
				return "", ErrKeyNotFound
			}
			return string(valueBytes), nil
		}
	}
	return "", ErrMissingPassword
}

// SlackConnector is an interface to the SlackConnector and is used to perform tests
type SlackConnector interface {
	GetAPIwithChannel(cm *ChatManager, chat *mysqlv1alpha1.Chat) (*slack.Client, string, error)
	PostMessage(api *slack.Client, channel string, message string) error
}

// DefaultSlackConnector an implementation of the SlackConnector
type DefaultSlackConnector struct {
}

// NewDefaultSlackConnector generates a Slack connector based on slack-go
func NewDefaultSlackConnector() SlackConnector {
	return &DefaultSlackConnector{}
}

// GetAPIwithChannel returns the API with the Channel ID for a named group or channel
func (cc *DefaultSlackConnector) GetAPIwithChannel(cm *ChatManager, chat *mysqlv1alpha1.Chat) (*slack.Client, string, error) {
	key := fmt.Sprintf("%s/%s", chat.Namespace, chat.Name)
	api, ok := cm.Reconciler.Chats[key]
	if !ok || api == nil {
		token, err := cm.getSlackToken(chat)
		if err != nil {
			cm.Reconciler.Log.Info(fmt.Sprintf("%s/%s connections to Slack failed with %v", chat.Namespace, chat.Name, err))
			return nil, "", ErrChatConnectionFailed
		}
		api = slack.New(token)
		cm.Reconciler.Chats[key] = api
	}
	next := ""
	for {
		conversation := &slack.GetConversationsParameters{
			Cursor:          next,
			ExcludeArchived: "true",
			Limit:           100,
			Types:           []string{"public_channel", "private_channel"},
		}
		channels, next, err := api.GetConversations(conversation)
		if err != nil {
			cm.Reconciler.Log.Info(fmt.Sprintf("%s/%s List conversations failed with %v", chat.Namespace, chat.Name, err))
			return nil, "", ErrChannelNotFound
		}
		for _, v := range channels {
			if v.Name == chat.Spec.Slack.Channel {
				return api, v.ID, nil
			}
		}
		if next == "" {
			cm.Reconciler.Log.Info(fmt.Sprintf("%s/%s Channel not found: %v", chat.Namespace, chat.Name, err))
			return nil, "", ErrChannelNotFound
		}
	}
}

// PostMessage Post a message on the API
func (cc *DefaultSlackConnector) PostMessage(api *slack.Client, channel string, message string) error {
	// TODO: check the api and channel are correct
	api.PostMessage(channel, slack.MsgOptionText(
		message,
		false,
	))
	return nil
}
