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

// ResourceImplementationSpec defines the desired state of ResourceImplementation
type ResourceImplementationSpec struct {
	Class          string                              `json:"class"`
	RequirementRef ImplementationRequirementDefinition `json:"requirementRef"`
}

type ImplementationRequirementDefinition struct {
	Code    string `json:"code"`
	Version string `json:"version"`
}

// ResourceImplementationStatus defines the observed state of ResourceImplementation
type ResourceImplementationStatus struct {
	ResourceAttestations []ResourceAttestation `json:"resourceAttestations"`
	TotalAttestations    int                   `json:"totalAttestations"`
	PassedAttestations   int                   `json:"passedAttestations"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ResourceImplementation is the Schema for the resourceimplementations API
type ResourceImplementation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceImplementationSpec   `json:"spec,omitempty"`
	Status ResourceImplementationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ResourceImplementationList contains a list of ResourceImplementation
type ResourceImplementationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceImplementation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceImplementation{}, &ResourceImplementationList{})
}
