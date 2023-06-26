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

// AttestationSpec defines the desired state of Attestation
type AttestationSpec struct {
	Type              AttestationType     `json:"type"`
	ImplementationRef string              `json:"implementationRef"`
	ProviderRef       AttestationProvider `json:"providerRef"`
}

type AttestationProvider struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type AttestationType struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// AttestationStatus defines the observed state of Attestation
type AttestationStatus struct {
	Childs []ResourceAttestationChild `json:"result"`
	Status string                     `json:"status"`
}

type ResourceAttestationChild struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Attestation is the Schema for the attestations API
type Attestation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AttestationSpec   `json:"spec,omitempty"`
	Status AttestationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AttestationList contains a list of Attestation
type AttestationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Attestation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Attestation{}, &AttestationList{})
}
