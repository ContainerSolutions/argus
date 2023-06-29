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

// ResourceRequirementSpec defines the desired state of ResourceRequirement
type ResourceRequirementSpec struct {
	Definition                    RequirementDefinition `json:"definition"`
	RequiredImplementationClasses []string              `json:"requiredImplementationClasses"`
}

// ResourceRequirementStatus defines the observed state of ResourceRequirement
type ResourceRequirementStatus struct {
	//+optional
	ApplicableResourceImplementations []NamespacedName `json:"applicableResourceImplementations,omitempty"`
	//+optional
	TotalImplementations int `json:"totalImplementations,omitempty"`
	//+optional
	ValidImplementations int `json:"validImplementations,omitempty"`
	//+optional
	Status string `json:"status,omitempty"`
	//+optional
	RequirementHash string `json:"requirementHash,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ResourceRequirement is the Schema for the resourcerequirements API
type ResourceRequirement struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceRequirementSpec   `json:"spec,omitempty"`
	Status ResourceRequirementStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ResourceRequirementList contains a list of ResourceRequirement
type ResourceRequirementList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceRequirement `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceRequirement{}, &ResourceRequirementList{})
}
