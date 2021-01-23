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
	// InstanceStoreNotReady the store is not ready
	InstanceStoreNotReady = "StoreNotReady"
	// InstanceStatefulSetInaccessible the statefulset cannot be accessed
	InstanceStatefulSetInaccessible = "StatefulSetInaccessible"
	// InstanceStatefulSetUpdated the statefulset has been updated
	InstanceStatefulSetUpdated = "StatefulSetUpdated"
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

// BackupScheduleSpec defines the backup schedule properties
type BackupScheduleSpec struct {
	// The backup store to use for backups
	Store string `json:"store,omitempty"`

	// The backup schedule to use for backups
	Schedule string `json:"schedule,omitempty"`
}

// MaintenanceScheduleSpec defines the backup schedule properties
type MaintenanceScheduleSpec struct {

	// The maintenance schedule
	Schedule string `json:"schedule,omitempty"`

	// The maintenance schedule in minutes
	Duration int `json:"duration,omitempty"`
}

// InstanceSpec defines the desired state of Instance
type InstanceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Restore when starting from an existing configuration
	Restore RestoreSpec `json:"restore,omitempty"`

	// Defines the backup schedules
	BackupSchedule BackupScheduleSpec `json:"backupSchedule,omitempty"`

	// Database is the default database name for the instance
	Database string `json:"database,omitempty"`

	// Defines the backup schedules
	MaintenanceSchedule MaintenanceScheduleSpec `json:"maintenanceSchedule,omitempty"`
}

// ScheduleEntry defines schedule properties
type ScheduleEntry struct {
	// The backup schedule that last applied
	Schedule string `json:"schedule,omitempty"`
	// The BackupJob ID in the Scheduler
	EntryID int `json:"entryID"`
}

// ScheduleStatus defines the schedule properties
type ScheduleStatus struct {
	// The Scheduler incarnation managed by the operator
	// +kubebuilder:default:="00000000-0000-0000-0000-000000000000"
	Incarnation string `json:"incarnation,omitempty"`
	// Properties for the backup schedule
	Backup ScheduleEntry `json:"backup,omitempty"`

	// Properties for the maintenance schedule
	Maintenance ScheduleEntry `json:"maintenance,omitempty"`

	// Properties to turn off the maintenance mode
	MaintenanceOff ScheduleEntry `json:"maintenanceOff,omitempty"`

	// When the maintenance mode should be turned off
	MaintenanceEndTime *metav1.Time `json:"maintenanceEndTime,omitempty"`
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
	// Defines if the database is currently in Maintenance Mode
	MaintenanceMode bool `json:"maintenanceMode"`
	// Conditions provides informations about the the last conditions
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// Schedules provides information about the current running schedules, including backups and maintenance
	Schedules ScheduleStatus `json:"schedules,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Instance ready"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.reason",description="Instance phase"
// +kubebuilder:printcolumn:name="Maintenance",type="boolean",JSONPath=".status.maintenanceMode",description="Instance currently in Maintenance"

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
