package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OperationMode is an Enum type to reference the mode of operation
type OperationMode string

const (
	// OperationModeImmediate defines operations that should apply immediately
	OperationModeImmediate OperationMode = "immediate"
	// OperationModeMaintenance defines operations that should apply on the next maintenance window
	OperationModeMaintenance OperationMode = "maintenance"
)

// OperationType is an Enum type to reference the type of operation
type OperationType string

const (
	// OperationTypeRestart a restart operation
	OperationTypeRestart OperationType = "restart"
)

// OperationSpec defines the desired state of Operation
type OperationSpec struct {

	// Mode defines how an operation is applied, it could be immediate or wait for the next maintenance window on the instance
	// +kubebuilder:validation:Enum=maintenance;immediate
	// +kubebuilder:default:="maintenance"
	Mode OperationMode `json:"mode,omitempty"`

	// Type defines the operation type. For now, only restart is supported
	// +kubebuilder:validation:Enum=restart
	// +kubebuilder:default:="restart"
	Type string `json:"type,omitempty"`

	// Defines the instance the operation applies to
	Instance string `json:"instance,omitempty"`
}

// OperationStatus defines the observed state of Operation
type OperationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Operation is the Schema for the operations API
type Operation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OperationSpec   `json:"spec,omitempty"`
	Status OperationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OperationList contains a list of Operation
type OperationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Operation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Operation{}, &OperationList{})
}
