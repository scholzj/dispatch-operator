package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RouterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Router `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Router struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              RouterSpec   `json:"spec"`
	Status            RouterStatus `json:"status,omitempty"`
}

type RouterSpec struct {
	// Nodes is the size of the memcached deployment
	Nodes int32 `json:"nodes,omitempty"`
}
type RouterStatus struct {
	// List of broker URLs
	URLs []string `json:"urls,omitempty"`
}
