package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags
// for the fields to be serialized.

const (
	// StoreCheckRequested specifies a state check request
	StoreCheckRequested = "CheckRequested"

	// StoreCheckSucceeded shows the last state check has passed
	StoreCheckSucceeded = "CheckSucceeded"

	// StoreCheckFailed shows the last state check has failed
	StoreCheckFailed = "CheckFailed"
)

// EnvVar represents an environment variable present in a store.
type EnvVar struct {
	// Name of the environment variable. Must be a C_IDENTIFIER.
	Name string `json:"name"`
	// Variable's value.
	// +optional
	Value string `json:"value,omitempty"`
	// Source for the environment variable's value. Cannot be used if value is
	// not empty.
	// +optional
	ValueFrom *EnvVarSource `json:"valueFrom,omitempty"`
}

// EnvVarSource represents a source for the value of an EnvVar.
type EnvVarSource struct {
	// Selects a key of a ConfigMap.
	// +optional
	ConfigMapKeyRef *corev1.ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`
	// Selects a key of a secret in the pod's namespace
	// +optional
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty"`
}

// StoreSpec defines the desired state of Store
type StoreSpec struct {
	// Defines the type of backend to be used for the store. For now on, only
	// s3 is supported (default: s3)
	// +kubebuilder:validation:Enum=s3
	// +kubebuilder:validation:Enum=blackhole
	Backend *string `json:"backend,omitempty"`
	// the store bucket
	Bucket string `json:"bucket"`
	// Prefix defines section of the path that will prefix files in the bucket.
	// This is to keep files from multiple sources in the same bucket.
	// +optional
	Prefix string `json:"prefix,omitempty"`
	// Envs defines a set of environment variables that can be used to access
	// secured stores which should be the case for every store
	// +optional
	Envs []corev1.EnvVar `json:"envs,omitempty"`
}

// StoreStatus defines the observed state of Store
type StoreStatus struct {
	// Defines if the store can be considered as ready or not
	Ready metav1.ConditionStatus `json:"ready,omitempty"`
	// Defines if the store current Reason
	Reason string `json:"reason,omitempty"`
	// A flag that indicates a resouce should be re-checked
	CheckRequested bool `json:"checkrequested,omitempty"`
	// A human readable message indicating details about why the store is in
	// this condition.
	Message string `json:"message,omitempty"`
	// A human readable message indicating details about why the store is in
	// this condition.
	Conditions []metav1.Condition `json:"Conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Store is the Schema for the stores API
type Store struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              StoreSpec   `json:"spec,omitempty"`
	Status            StoreStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StoreList contains a list of Store
type StoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Store `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Store{}, &StoreList{})
}
