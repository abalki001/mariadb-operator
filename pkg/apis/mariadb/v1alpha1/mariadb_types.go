package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MariaDBSpec defines the desired state of MariaDB
type MariaDBSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Size is the size of the deployment
	Size int32 `json:"size"`

	// Database additional user details (base64 encoded)
	Username string `json:"username"`

	// Database additional user password (base64 encoded)
	Password string `json:"password"`

	// New Database name
	Database string `json:"database"`

	// Root user password
	Rootpwd string `json:"rootpwd"`

	// Image name with version
	Image string `json:"image"`

	// Database storage Path
	DataStoragePath string `json:"dataStoragePath"`

	// Database storage Size (Ex. 1Gi, 100Mi)
	DataStorageSize string `json:"dataStorageSize"`

	// Port number exposed for Database service
	Port int32 `json:"port"`
}

// MariaDBStatus defines the observed state of MariaDB
type MariaDBStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Nodes are the names of the pods
	Nodes []string `json:"nodes,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MariaDB is the Schema for the mariadbs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=mariadbs,scope=Namespaced
type MariaDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MariaDBSpec   `json:"spec,omitempty"`
	Status MariaDBStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MariaDBList contains a list of MariaDB
type MariaDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MariaDB `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MariaDB{}, &MariaDBList{})
}
