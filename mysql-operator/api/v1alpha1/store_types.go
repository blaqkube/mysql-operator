package v1alpha1

import (
	"github.com/operator-framework/operator-lib/status"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EnvVar represents an environment variable present in a store.
type EnvVar struct {
	// Name of the environment variable. Must be a C_IDENTIFIER.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	// Variable's value.
	// +optional
	Value string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value"`
	// Source for the environment variable's value. Cannot be used if value is not empty.
	// +optional
	ValueFrom *EnvVarSource `json:"valueFrom,omitempty" protobuf:"bytes,3,opt,name=valueFrom"`
}

// EnvVarSource represents a source for the value of an EnvVar.
type EnvVarSource struct {
	// Selects a key of a ConfigMap.
	// +optional
	ConfigMapKeyRef *corev1.ConfigMapKeySelector `json:"configMapKeyRef,omitempty" protobuf:"bytes,3,opt,name=configMapKeyRef"`
	// Selects a key of a secret in the pod's namespace
	// +optional
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty" protobuf:"bytes,4,opt,name=secretKeyRef"`
}

// S3Info defines store informations for S3
type S3Info struct {
	Bucket string          `json:"bucket"`
	Path   string          `json:"path"`
	Env    []corev1.EnvVar `json:"env,omitempty"`
}

// StoreSpec defines the desired state of Store
type StoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Enum=s3
	Backend *string `json:"backend,omitempty"`

	// S3 provides the information to backup databases on S3
	S3 *S3Info `json:"s3,omitempty"`
}

// StoreStatus defines the observed state of Store
type StoreStatus struct {

	// A brief CamelCase message indicating details about why the store is in this state. e.g. 'Ready'
	Reason string `json:"reason,omitempty"`

	// A human readable message indicating details about why the store is in this condition.
	Message string `json:"message,omitempty"`

	// Conditions that provides an history of the store conditions
	Conditions []status.Conditions `json:"store_conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Store is the Schema for the stores API
type Store struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StoreSpec   `json:"spec,omitempty"`
	Status StoreStatus `json:"status,omitempty"`
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
