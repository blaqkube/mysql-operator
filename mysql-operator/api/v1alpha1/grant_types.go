package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	// GrantUserDatabaseMismatch the associated user and database do not match the same instance
	GrantUserDatabaseMismatch = "UserDatabaseMismatch"
	// GrantUserAccessError the associated user could not be accessed
	GrantUserAccessError = "UserAccessError"
	// GrantUserNotReady the associated user is not yet ready
	GrantUserNotReady = "UserNotReady"
	// GrantDatabaseAccessError the associated database could not be accessed
	GrantDatabaseAccessError = "DatabaseAccessError"
	// GrantDatabaseNotReady the associated database is not yet ready
	GrantDatabaseNotReady = "DatabaseNotReady"
	// GrantInstanceAccessError the associated instance could not be accessed
	GrantInstanceAccessError = "InstanceAccessError"
	// GrantInstanceNotReady the associated instance is not yet ready
	GrantInstanceNotReady = "InstanceNotReady"
	// GrantAgentNotFound the agent could not be found
	GrantAgentNotFound = "AgentNotFound"
	// GrantAgentFailed a request to the agent failed
	GrantAgentFailed = "AgentFailed"
	// GrantNotImplemented grant creation has not been implemented
	GrantNotImplemented = "NotImplemented"
	// GrantSucceeded grant creation has succeeded
	GrantSucceeded = "Succeeded"
)

// AccessMode is an Enum type to reference different storages
type AccessMode string

const (
	// AccessReadWrite Read and write access for a user on a database
	AccessReadWrite AccessMode = "readWrite"
	// AccessReadOnly Read-only access for a user on a database
	AccessReadOnly AccessMode = "readOnly"
)

// GrantSpec defines the desired state of Grant
type GrantSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Defines the granted user
	User string `json:"user"`
	// Defines the granted database
	Database string `json:"database"`
	// Defines the type of access for the user and database
	// +kubebuilder:validation:Enum=readWrite;readOnly
	AccessMode string `json:"accessMode"`
}

// GrantStatus defines the observed state of Grant
type GrantStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
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
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Grant ready"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.Reason",description="Grant phase"

// Grant is the Schema for the grants API
type Grant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GrantSpec   `json:"spec,omitempty"`
	Status GrantStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GrantList contains a list of Grant
type GrantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Grant `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Grant{}, &GrantList{})
}
