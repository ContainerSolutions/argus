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

// ComponentSpec defines the desired state of Component
type ComponentSpec struct {
	Type    string   `json:"type"`
	Classes []string `json:"classes"`
	Parents []string `json:"parents"`
}

// ComponentStatus defines the observed state of Component
type ComponentStatus struct {
	//+kubebuilder:default=0
	TotalControls int `json:"totalControls"`
	//+kubebuilder:default=0
	ImplementedControls int `json:"implementedControls"`
	//+optional
	Children map[string]ComponentChild `json:"children,omitempty"`
	//+optional
	Controls map[string]*ComponentControlCompliance `json:"Controls,omitempty"`
	//+kubebuilder:default=0
	TotalChildren int `json:"totalChildren"`
	//+kubebuilder:default=0
	CompliantChildren int `json:"compliantChildren"`
	//+optional
	RunAt metav1.Time `json:"runAt,omitempty"`
}

type ComponentControlCompliance struct {
	Implemented bool `json:"implemented"`
}

// All parent relationship is flattened. TODO - maybe we want to have the whole hierarchy here?
// TODO - If the child has a Control the parent does not have (and it is non compliant to that Control)
// Should the parent be marked as non compliant? Or rather just as having Non compliant Children?
// TODO - Need a way to check compliance based on Control Classes
type ComponentChild struct {
	Compliant bool `json:"compliant"`
}

// Component is the Schema for the Components API

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Total Controls",type=integer,JSONPath=`.status.totalControls`
// +kubebuilder:printcolumn:name="Implemented Controls",type=integer,JSONPath=`.status.implementedControls`
// +kubebuilder:printcolumn:name="Last Run",type=string,JSONPath=`.status.runAt`
type Component struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComponentSpec   `json:"spec,omitempty"`
	Status ComponentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ComponentList contains a list of Component
type ComponentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Component `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Component{}, &ComponentList{})
}
