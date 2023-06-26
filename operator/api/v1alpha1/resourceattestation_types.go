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

// ResourceAttestationSpec defines the desired state of ResourceAttestation
type ResourceAttestationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ResourceAttestation. Edit resourceattestation_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// ResourceAttestationStatus defines the observed state of ResourceAttestation
type ResourceAttestationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ResourceAttestation is the Schema for the resourceattestations API
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
