package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BackupSpec defines the desired state of Backup
type BackupSpec struct {
	// The store to use to perform the backup.
	Store string `json:"store"`
	// Instance to backup.
	Instance string `json:"instance"`
}

const (
	// BackupFailed the associated backup has failed
	BackupFailed = "Failed"
	// BackupRunning the associated backup is running
	BackupRunning = "Running"
	// BackupStoreAccessError the associated store could not be accessed
	BackupStoreAccessError = "StoreAccessError"
	// BackupStoreNotReady the associated instance is not yet ready
	BackupStoreNotReady = "StoreNotReady"
	// BackupMissingVariable some variables are missing
	BackupMissingVariable = "StoreMissingVariable"
	// BackupInstanceAccessError the associated instance could not be accessed
	BackupInstanceAccessError = "InstanceAccessError"
	// BackupInstanceNotReady the associated instance is not yet ready
	BackupInstanceNotReady = "InstanceNotReady"
	// BackupAgentNotFound the agent could not be found
	BackupAgentNotFound = "AgentNotFound"
	// BackupAgentFailed a request to the agent failed
	BackupAgentFailed = "AgentFailed"
	// BackupNotImplemented grant creation has not been implemented
	BackupNotImplemented = "NotImplemented"
	// BackupSucceeded grant creation has succeeded
	BackupSucceeded = "Succeeded"
)

// BackupDetails defines the Backup Location and StartupTime
type BackupDetails struct {
	// Internal Identifier
	Identifier string `json:"identifier,omitempty"`
	// Bucket
	Bucket string `json:"bucket,omitempty"`
	// Location in bucket
	Location string `json:"location,omitempty"`
	// Start Time
	StartTime *metav1.Time `json:"backupTime,omitempty"`
	// End Time
	EndTime *metav1.Time `json:"endTime,omitempty"`
}

// BackupStatus defines the observed state of Backup
type BackupStatus struct {
	// Defines the details for the backup
	Details *BackupDetails `json:"details,omitempty"`
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

// Backup is the Schema for the backups API
type Backup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupSpec   `json:"spec,omitempty"`
	Status BackupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BackupList contains a list of Backup
type BackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Backup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Backup{}, &BackupList{})
}
