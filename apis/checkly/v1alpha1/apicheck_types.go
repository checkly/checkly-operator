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

	// Locations determines the locations where the checks are run from, see https://www.checklyhq.com/docs/monitoring/global-locations/ for a list, use AWS Region codes, ex. eu-west-1 for Ireland
	Locations []string `json:"locations,omitempty"`

	// Endpoint determines which URL to monitor, ex. https://foo.bar/baz
	Endpoint string `json:"endpoint"`

	// Success determines the returned success code, ex. 200
	Success string `json:"success"`

	// MaxResponseTime determines what the maximum number of miliseconds can pass before the check fails, default 15000
	MaxResponseTime int `json:"maxresponsetime,omitempty"`

	// Group determines in which group does the check belong to
	Group string `json:"group"`
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
