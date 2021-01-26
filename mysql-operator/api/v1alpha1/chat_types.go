package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// ChatPending did not attempt to connect yet or did not get result
	ChatPending = "Pending"
	// ChatSlackConnectionError could not connect to Slack
	ChatSlackConnectionError = "SlackConnectionError"
	// ChatSlackChannelError could not connect to the Slack Channel
	ChatSlackChannelError = "SlackChannelError"
	// ChatSucceeded could connect to Slack and send a message
	ChatSucceeded = "Succeeded"
)

// SlackSpec defines the properties associated with Slack
type SlackSpec struct {

	// Token's value
	// +optional
	Token string `json:"token,omitempty"`
	// Source for the Token's value. Cannot be used if Token is
	// not empty.
	// +optional
	TokenFrom *TokenSource `json:"tokenFrom,omitempty"`

	// Channel is the name of the channel to use
	Channel string `json:"channel,omitempty"`
}

// TokenSource represents a source for the value of a Token.
type TokenSource struct {
	// Selects a key of a ConfigMap.
	// +optional
	ConfigMapKeyRef *corev1.ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`
	// Selects a key of a secret in the pod's namespace
	// +optional
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty"`
}

// ChatSpec defines the desired state of Chat
type ChatSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Slack provides a slack configuration
	Slack SlackSpec `json:"slack,omitempty"`
}

// ChatStatus defines the observed state of Chat
type ChatStatus struct {
	// Defines if slack can be considered as ready or not
	Ready metav1.ConditionStatus `json:"ready,omitempty"`
	// Defines the Reason behind Slack can be considered as ready or not
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about why the chant is in
	// this condition.
	Message string `json:"message,omitempty"`
	// Allow to understand the history of conditions
	Conditions []metav1.Condition `json:"Conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Chat ready"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.reason",description="Chat phase"

// Chat is the Schema for the chats API
type Chat struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChatSpec   `json:"spec,omitempty"`
	Status ChatStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ChatList contains a list of Chat
type ChatList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Chat `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Chat{}, &ChatList{})
}
