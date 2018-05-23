package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConsulObject describes a Consul Object of KVs.
type ConsulObject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ConsulObjectSpec `json:"spec"`
}

// ConsulObjectSpec is the spec for a Consul Object resource
type ConsulObjectSpec struct {
	KV *[]ConsulObjectKV `json:"kv"`
}

type ConsulObjectKV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Flag  string `json:"flag,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConsulObjectList is a list of ConsulObjects resources
type ConsulObjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ConsulObject `json:"items"`
}
