/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GroupSpec defines the desired state of Group
type GroupSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Locations determines the locations where the checks are run from, see https://www.checklyhq.com/docs/monitoring/global-locations/ for a list, use AWS Region codes, ex. eu-west-1 for Ireland
	Locations []string `json:"locations,omitempty"`

	// Locations determines the locations where the checks are run from, see https://www.checklyhq.com/docs/monitoring/global-locations/ for a list, use AWS Region codes, ex. eu-west-1 for Ireland
	PrivateLocations []string `json:"privateLocations,omitempty"`

	// Activated determines if the created group is muted or not, default false
	Activated bool `json:"muted,omitempty"`

	// AlertChannels determines where to send alerts
	AlertChannels []string `json:"alertchannel,omitempty"`
}

// GroupStatus defines the observed state of Group
type GroupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ID holds the ID of the created checklyhq.com group
	ID int64 `json:"ID"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// Group is the Schema for the groups API
type Group struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GroupSpec   `json:"spec,omitempty"`
	Status GroupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GroupList contains a list of Group
type GroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Group `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Group{}, &GroupList{})
}
