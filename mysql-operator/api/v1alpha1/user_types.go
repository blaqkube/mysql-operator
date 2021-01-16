package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	// UserPasswordError the password was not defined properly
	UserPasswordError = "PasswordError"
	// UserPasswordAccessError the password could not be access
	UserPasswordAccessError = "PasswordAccessError"
	// UserInstanceAccessError the associated instance could not be accessed
	UserInstanceAccessError = "InstanceAccessError"
	// UserInstanceNotReady the associated instance is not yet ready
	UserInstanceNotReady = "InstanceNotReady"
	// UserAgentNotFound the agent could not be found
	UserAgentNotFound = "AgentNotFound"
	// UserAgentFailed a request to the agent failed
	UserAgentFailed = "AgentFailed"
	// UserSucceeded user creation has succeeded
	UserSucceeded = "Succeeded"
)

// UserSpec defines the desired state of User
type UserSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Instance string `json:"instance"`
	Username string `json:"username"`
	// Password's value.
	// +optional
	Password string `json:"password,omitempty"`
	// Source for the environment Password's value. Cannot be used if Password is
	// not empty.
	// +optional
	PasswordFrom *PasswordSource `json:"passwordFrom,omitempty"`
}

// PasswordSource represents a source for the value of a Password.
type PasswordSource struct {
	// Selects a key of a ConfigMap.
	// +optional
	ConfigMapKeyRef *corev1.ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`
	// Selects a key of a secret in the pod's namespace
	// +optional
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty"`
}

// UserStatus defines the observed state of User
type UserStatus struct {
	// Defines if the store can be considered as ready or not
	Ready metav1.ConditionStatus `json:"ready,omitempty"`
	// Defines if the store current Reason
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about why the store is in
	// this condition.
	Message string `json:"message,omitempty"`
	// A human readable message indicating details about why the store is in
	// this condition.
	Conditions []metav1.Condition `json:"Conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="User ready"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.Reason",description="User phase"

// User is the Schema for the users API
type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UserSpec   `json:"spec,omitempty"`
	Status UserStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// UserList contains a list of User
type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []User `json:"items"`
}

func init() {
	SchemeBuilder.Register(&User{}, &UserList{})
}
