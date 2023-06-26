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

// AttestationSpec defines the desired state of Attestation
type AttestationSpec struct {
	Type              AttestationType `json:"type"`
	ImplementationRef string          `json:"implementationRef"`
}

type AttestationType struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// AttestationStatus defines the observed state of Attestation
type AttestationStatus struct {
	Result map[string]AttestationResult `json:"result"`
	Status string                       `json:"status"`
}

type AttestationResult struct {
	Logs   string
	Result AttestationResultType
	Reason string
	Err    string
	RunAt  metav1.Time
}
type AttestationResultType string

const (
	AttestationResultTypePass       AttestationResultType = "Pass"
	AttestationResultTypeFail       AttestationResultType = "Fail"
	AttestationResultTypeUnknown    AttestationResultType = "Unknown"
	AttestationResultTypeNotStarted AttestationResultType = "Not Started"
)

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
