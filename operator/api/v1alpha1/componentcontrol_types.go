/*
Copyright 2023.

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

// ComponentControlSpec defines the desired state of ComponentControl
type ComponentControlSpec struct {
	Definition                ControlDefinition `json:"definition"`
	RequiredAssessmentClasses []string          `json:"requiredAssessmentClasses"`
}

// ComponentControlStatus defines the observed state of ComponentControl
type ComponentControlStatus struct {
	//+optional
	ApplicableComponentAssessments []NamespacedName `json:"applicableComponentAssessments,omitempty"`
	//+kubebuilder:default=0
	TotalAssessments int `json:"totalAssessments"`
	//+kubebuilder:default=0
	ValidAssessments int `json:"validAssessments"`
	//+optional
	Status string `json:"status,omitempty"`
	//+optional
	ControlHash string `json:"ControlHash,omitempty"`
	//+optional
	RunAt metav1.Time `json:"runAt,omitempty"`
}

// ComponentControl is the Schema for the ComponentControls API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Total Assessments",type=integer,JSONPath=`.status.totalAssessments`
// +kubebuilder:printcolumn:name="Valid Assessments",type=integer,JSONPath=`.status.validAssessments`
// +kubebuilder:printcolumn:name="Last Run",type=string,JSONPath=`.status.runAt`
type ComponentControl struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComponentControlSpec   `json:"spec,omitempty"`
	Status ComponentControlStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ComponentControlList contains a list of ComponentControl
type ComponentControlList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ComponentControl `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ComponentControl{}, &ComponentControlList{})
}
