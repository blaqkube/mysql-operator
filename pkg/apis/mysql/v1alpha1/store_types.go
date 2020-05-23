package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// S3Credentials defines the credentials to connect to S3
type S3Credentials struct {
	AccessKey string `json:"aws_access_key_id"`
	SecretKey string `json:"aws_secret_access_key"`
	Region string `json:"region"`
}

// S3Info defines the S3 bucket, path and connection
type S3Info struct {
	Bucket string `json:"bucket"`
	Path string `json:"path,omitempty"`
	Credentials S3Credentials `json:"credentials,omitempty"`
}

// StoreSpec defines the desired state of Store
type StoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// +kubebuilder:validation:Enum=s3
	Backend string `json:"backend"`
	S3Access S3Info `json:"s3access,omitempty"`
}

// StoreStatus defines the observed state of Store
type StoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	LastConnection string `json:"last_connection"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Store is the Schema for the stores API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=stores,scope=Namespaced
type Store struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StoreSpec   `json:"spec,omitempty"`
	Status StoreStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StoreList contains a list of Store
type StoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Store `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Store{}, &StoreList{})
}
