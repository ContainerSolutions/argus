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

// ImplementationSpec defines the desired state of Implementation
type ImplementationSpec struct {
	Class string `json:"class"`
	//+default="Cascade"
	CascadePolicy  ImplementationCascadePolicy         `json:"cascadePolicy"`
	RequirementRef ImplementationRequirementDefinition `json:"requirementRef"`
	ResourceRef    []NamespacedName                    `json:"resourceRef"`
}

type ImplementationCascadePolicy string

const (
	CascadingPolicyCascade ImplementationCascadePolicy = "Cascade"
	CascadingPolicyNone    ImplementationCascadePolicy = "None"
)

// ImplementationStatus defines the observed state of Implementation
type ImplementationStatus struct {
	Childs []NamespacedName `json:"childs"`
	Status string           `json:"status"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Implementation is the Schema for the implementations API
type Implementation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImplementationSpec   `json:"spec,omitempty"`
	Status ImplementationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ImplementationList contains a list of Implementation
type ImplementationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Implementation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Implementation{}, &ImplementationList{})
}
