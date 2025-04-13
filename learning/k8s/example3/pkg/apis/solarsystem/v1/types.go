package v1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Star struct {
	v1.TypeMeta   `json:",inline"`
	v1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StarSpec
	Status StarStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type StarList struct {
	v1.TypeMeta `json:",inline"`
	v1.ListMeta `json:"metadata"`

	Items []Star `json:"items"`
}

type StarSpec struct {
	Type     string `json:"type"`
	Location string `json:"location"`
}

type StarStatus struct {
	Status string `json:",inline"`
}
