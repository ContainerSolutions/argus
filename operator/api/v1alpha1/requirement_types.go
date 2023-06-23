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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RequirementSpec defines the desired state of Requirement
type RequirementSpec struct {
	Version  string `json:"version"`
	Code     string `json:"code"`
	Class    string `json:"class"`
	Category string `json:"category"`
	// TODO define classes objects instead of a free string? Strong typing means better validation
	ApplicableResourceClasses     []string `json:"applicableResourceClasses"`
	RequiredImplementationClasses []string `json:"requiredImplementationClasses"`
}

// RequirementStatus defines the observed state of Requirement
type RequirementStatus struct {
	RequirementHash string `json:"requirementHash"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Requirement is the Schema for the requirements API
type Requirement struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RequirementSpec   `json:"spec,omitempty"`
	Status RequirementStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RequirementList contains a list of Requirement
type RequirementList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Requirement `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Requirement{}, &RequirementList{})
}