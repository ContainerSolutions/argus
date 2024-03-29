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

// ControlSpec defines the desired state of Control
type ControlSpec struct {
	Definition ControlDefinition `json:"definition"`
	// TODO define classes objects instead of a free string? Strong typing means better validation
	ApplicableComponentClasses []string `json:"applicableComponentClasses"`
	RequiredAssessmentClasses  []string `json:"requiredAssessmentClasses"`
}

type ControlDefinition struct {
	Version     string `json:"version"`
	Code        string `json:"code"`
	Class       string `json:"class"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

// ControlStatus defines the observed state of Control
type ControlStatus struct {
	//+optional
	Children    []NamespacedName `json:"children,omitempty"`
	ControlHash string           `json:"ControlHash"`
}

type NamespacedName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Control is the Schema for the Controls API
type Control struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ControlSpec   `json:"spec,omitempty"`
	Status ControlStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ControlList contains a list of Control
type ControlList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Control `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Control{}, &ControlList{})
}
