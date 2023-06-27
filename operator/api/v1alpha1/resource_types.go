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

// ResourceSpec defines the desired state of Resource
type ResourceSpec struct {
	Type    string   `json:"type"`
	Classes []string `json:"classes"`
	Parents []string `json:"parents"`
}

// ResourceStatus defines the observed state of Resource
type ResourceStatus struct {
	//+optional
	TotalRequirements int `json:"totalRequirements,omitempty"`
	//+optional
	ImplementedRequirements int `json:"implementedRequirements,omitempty"`
	//+optional
	Children map[string]ResourceChild `json:"children,omitempty"`
	//+optional
	Requirements map[string]*ResourceRequirementCompliance `json:"requirements,omitempty"`
	//+optional
	TotalChildren int `json:"totalChildren,omitempty"`
	//+optional
	CompliantChildren int `json:"compliantChildren,omitempty"`
}

type ResourceRequirementCompliance struct {
	Implemented bool `json:"implemented"`
}

// All parent relationship is flattened. TODO - maybe we want to have the whole hierarchy here?
// TODO - If the child has a requirement the parent does not have (and it is non compliant to that requirement)
// Should the parent be marked as non compliant? Or rather just as having Non compliant Children?
// TODO - Need a way to check compliance based on requirement Classes
type ResourceChild struct {
	Compliant bool `json:"compliant"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Resource is the Schema for the resources API
type Resource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceSpec   `json:"spec,omitempty"`
	Status ResourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ResourceList contains a list of Resource
type ResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Resource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Resource{}, &ResourceList{})
}
