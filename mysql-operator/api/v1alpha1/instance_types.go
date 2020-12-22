package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	// InstanceInitializing instance creation has been requested
	InstanceInitializing = "Initializing"
	// InstanceExporterSecretInaccessible the secret for the exporter could not be accessed
	InstanceExporterSecretInaccessible = "ExporterSecretInaccessible"
	// InstanceExporterSecretCreated the secret for the exporter has been created
	InstanceExporterSecretCreated = "ExporterSecretCreated"
	// InstanceExporterSecretDeleted the secret for the exporter has been deleted
	InstanceExporterSecretDeleted = "ExporterSecretDeleted"
	// InstanceExporterSecretFailed the secret for the exporter could not be created
	InstanceExporterSecretFailed = "ExporterSecretFailed"
	// InstanceStoreInaccessible the store cannot be accessed
	InstanceStoreInaccessible = "StoreInaccessible"
	// InstanceStatefulSetInaccessible the store cannot be accessed
	InstanceStatefulSetInaccessible = "StatefulSetInaccessible"
	// InstanceStatefulSetFailed the statefulset could not be created
	InstanceStatefulSetFailed = "StatefulSetFailed"
	// InstanceStatefulSetCreated the statefulset has been successfully created
	InstanceStatefulSetCreated = "StatefulSetCreated"
	// InstanceStatefulSetWaiting the statefulset is not yet reported as ready
	InstanceStatefulSetWaiting = "StatefulSetWaitingForReady"
	// InstanceStatefulSetReady the statefulset is ready
	InstanceStatefulSetReady = "StatefulSetReady"
)

// RestoreSpec defines the backup location when create a instance with a restore
type RestoreSpec struct {
	Store string `json:"store,omitempty"`

	Location string `json:"location,omitempty"`
}

// InstanceSpec defines the desired state of Instance
type InstanceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Restore when starting from an existing configuration
	Restore RestoreSpec `json:"restore,omitempty"`

	// Database is the default database name for the instance
	Database string `json:"database,omitempty"`
}

// InstanceStatus defines the observed state of Instance
type InstanceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// StatefulSet keeps track of the instance Statefulset
	StatefulSet corev1.ObjectReference `json:"statefulset,omitempty"`
	// ExporterSecret keeps track of the secret used for the exporter
	ExporterSecret corev1.ObjectReference `json:"exporter,omitempty"`
	// Defines if the instance is ready
	Ready metav1.ConditionStatus `json:"ready,omitempty"`
	// Defines if the store current Reason
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about why the store is in
	// this condition.
	Message string `json:"message,omitempty"`
	// Conditions provides informations about the the last conditions
	Conditions []metav1.Condition `json:"Conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Instance is the Schema for the instances API
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstanceSpec   `json:"spec,omitempty"`
	Status InstanceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// InstanceList contains a list of Instance
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Instance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Instance{}, &InstanceList{})
}
