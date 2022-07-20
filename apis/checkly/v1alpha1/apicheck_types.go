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

// ApiCheckSpec defines the desired state of ApiCheck
type ApiCheckSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Frequency is used to determine the frequency of the checks in minutes, default 5
	Frequency int `json:"frequency,omitempty"`

	// Muted determines if the created alert is muted or not, default false
	Muted bool `json:"muted,omitempty"`

	// Endpoint determines which URL to monitor, ex. https://foo.bar/baz
	Endpoint string `json:"endpoint"`

	// Success determines the returned success code, ex. 200
	Success string `json:"success"`

	// MaxResponseTime determines what the maximum number of miliseconds can pass before the check fails, default 15000
	MaxResponseTime int `json:"maxresponsetime,omitempty"`

	// Group determines in which group does the check belong to
	Group string `json:"group"`

	// GroupNamespace determine in which namespace was the group defined
	GroupNamespace string `json:"groupnamespace,omitempty"`
}

// ApiCheckStatus defines the observed state of ApiCheck
type ApiCheckStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ID holds the checklyhq.com internal ID of the check
	ID string `json:"id"`

	// GroupID holds the ID of the group where the check belongs to
	GroupID int64 `json:"groupId"`
}

//+kubebuilder:object:root=true
//+kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".spec.endpoint",description="Name of the monitored endpoint"
//+kubebuilder:printcolumn:name="Status code",type="string",JSONPath=".spec.success",description="Expected status code"
//+kubebuilder:printcolumn:name="Muted",type="boolean",JSONPath=".spec.muted"
//+kubebuilder:printcolumn:name="Group",type="string",JSONPath=".spec.group"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
//+kubebuilder:subresource:status

// ApiCheck is the Schema for the apichecks API
type ApiCheck struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiCheckSpec   `json:"spec,omitempty"`
	Status ApiCheckStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApiCheckList contains a list of ApiCheck
type ApiCheckList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApiCheck `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApiCheck{}, &ApiCheckList{})
}
