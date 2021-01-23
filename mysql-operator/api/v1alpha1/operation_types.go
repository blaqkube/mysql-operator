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

const (
	// OperationWaitingForMaintenanceWindow operation is waiting for maintenance window
	OperationWaitingForMaintenanceWindow = "WaitingMaintenanceWindow"
	// OperationRunning the operation is currently running
	OperationRunning = "Running"
	// OperationExecutedWithFailure the operation is executed and it has failed
	OperationExecutedWithFailure = "ExecutedWithFailure"
	// OperationExecutedWithSuccess the operation is executed with success
	OperationExecutedWithSuccess = "ExecutedWithSuccess"
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
	// Defines the current Reason for the operation
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the operation and the
	// associated condition.
	Message string `json:"message,omitempty"`
	// Conditions provides informations about the the last conditions
	Conditions []metav1.Condition `json:"conditions,omitempty"`
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
