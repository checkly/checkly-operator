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
	"github.com/checkly/checkly-go-sdk"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AlertChannelSpec defines the desired state of AlertChannel
type AlertChannelSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// SendRecovery determines if the Recovery event should be sent to the alert channel
	SendRecovery bool `json:"sendrecovery,omitempty"`

	// SendFailure determines if the Failure event should be sent to the alerting channel
	SendFailure bool `json:"sendfailure,omitempty"`

	// OpsGenie holds information about the Opsgenie alert configuration
	OpsGenie AlertChannelOpsGenie `json:"opsgenie,omitempty"`

	// Email holds information about the Email alert configuration
	Email checkly.AlertChannelEmail `json:"email,omitempty"`
}

type AlertChannelOpsGenie struct {
	// APISecret determines where the secret ref is to pull the OpsGenie API key from
	APISecret corev1.ObjectReference `json:"apisecret"`

	// Region holds information about the OpsGenie region (EU or US)
	Region string `json:"region,omitempty"`

	// Priority assigned to the alerts sent from checklyhq.com
	Priority string `json:"priority,omitempty"`
}

// AlertChannelStatus defines the observed state of AlertChannel
type AlertChannelStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ID int64 `json:"id"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// AlertChannel is the Schema for the alertchannels API
type AlertChannel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlertChannelSpec   `json:"spec,omitempty"`
	Status AlertChannelStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AlertChannelList contains a list of AlertChannel
type AlertChannelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlertChannel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AlertChannel{}, &AlertChannelList{})
}
