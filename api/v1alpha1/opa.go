package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// OpaPolicySpec defines the desired state of the OPAPolicy
type OpaPolicySpec struct {
	DeploymentName string `json:"deploymentName"`

	PolicyName string `json:"policyName"`

	Policy string `json:"policy"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// OpaPolicy is the Schema for the OpaPolicy
type OpaPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec OpaPolicySpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// OpaPolicyList contains a list of OpaPolicies
type OpaPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []OpaPolicy `json:"items"`
}

// OpaMessage is used by the watcher to inform the client of OPA policy changes
type OpaMessage struct {
	Method    string
	OpaPolicy OpaPolicy
}
