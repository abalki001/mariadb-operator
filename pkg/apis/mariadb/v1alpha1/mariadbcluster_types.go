package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MariaDBClusterSpec defines the desired state of MariaDBCluster
type MariaDBClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

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

	// Cluster Configuration
	Cluster ClusterDefinitionStruct `json:"cluster,omitempty"`
}

// Cluster definitions
type ClusterDefinitionStruct struct {
	Enabled bool `json:"enabled,omitempty"`

	// Specifies if this is first node in cluster and will initiate a new cluster
	FirstNode bool `json:"firstNode,omitempty"`

	// Name of Node where this pod instance is to be deployed
	NodeName string `json:"nodeName"`
}

// MariaDBClusterStatus defines the observed state of MariaDBCluster
type MariaDBClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MariaDBCluster is the Schema for the mariadbclusters API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=mariadbclusters,scope=Namespaced
type MariaDBCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MariaDBClusterSpec   `json:"spec,omitempty"`
	Status MariaDBClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MariaDBClusterList contains a list of MariaDBCluster
type MariaDBClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MariaDBCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MariaDBCluster{}, &MariaDBClusterList{})
}
