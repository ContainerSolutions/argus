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

// ResourceAttestationSpec defines the desired state of ResourceAttestation
type ResourceAttestationSpec struct {
	ProviderRef AttestationProviderRef `json:"providerRef"`
}

// ResourceAttestationStatus defines the observed state of ResourceAttestation
type ResourceAttestationStatus struct {
	Result AttestationResult `json:"result"`
	Status string            `json:"status"`
}

type AttestationResult struct {
	Logs   string                `json:"logs"`
	Result AttestationResultType `json:"result"`
	Reason string                `json:"reason"`
	//+optional
	Err   string      `json:"err,omitempty"`
	RunAt metav1.Time `json:"runAt"`
}
type AttestationResultType string

const (
	AttestationResultTypePass       AttestationResultType = "Pass"
	AttestationResultTypeFail       AttestationResultType = "Fail"
	AttestationResultTypeUnknown    AttestationResultType = "Unknown"
	AttestationResultTypeNotStarted AttestationResultType = "Not Started"
)

// ResourceAttestation is the Schema for the resourceattestations API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Result",type=string,JSONPath=`.status.result.result`
// +kubebuilder:printcolumn:name="Last Run",type=string,JSONPath=`.status.result.runAt`
type ResourceAttestation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceAttestationSpec   `json:"spec,omitempty"`
	Status ResourceAttestationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ResourceAttestationList contains a list of ResourceAttestation
type ResourceAttestationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceAttestation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceAttestation{}, &ResourceAttestationList{})
}
