package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MaintenanceSpec defines the maintenance window for operations
type MaintenanceSpec struct {
	WindowStart string `json:"windowStart,omitempty"`
	Backup      bool   `json:"backup,omitempty"`
	BackupStore string `json:"backupStore,omitempty"`
}

// InstanceSpec defines the desired state of Instance
// RestoreSpec defines the backup location when create a instance with a restore
type RestoreSpec struct {
	Store string `json:"store,omitempty"`

	FilePath string `json:"filePath,omitempty"`
}

// InstanceSpec defines the desired state of Instance
type InstanceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// New Database name
	Database    string          `json:"database"`
	Maintenance MaintenanceSpec `json:"maintenance,omitempty"`
	// Restore when starting from an existing configuration
	Restore RestoreSpec `json:"restore,omitempty"`
	Version string      `json:"version,omitempty"`
}

// InstanceStatus defines the observed state of Instance
type InstanceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// Nodes are the names of the pods
	Node          string            `json:"node,omitempty"`
	LastCondition string            `json:"lastCondition,omitempty"`
	Conditions    []ConditionStatus `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Instance is the Schema for the instances API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=instances,scope=Namespaced
// +kubebuilder:storageversion
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstanceSpec   `json:"spec,omitempty"`
	Status InstanceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstanceList contains a list of Instance
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Instance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Instance{}, &InstanceList{})
}
