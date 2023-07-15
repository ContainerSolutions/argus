package attestation

import (
	"context"
	"testing"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestGetComponentAttestations(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name           string
		attestation    *argusiov1alpha1.Attestation
		expectedOutput map[string]argusiov1alpha1.ComponentAttestation
		expectedError  string
		cl             client.Client
	}{
		{
			name:           "no op",
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestation:    makeAttestation(),
			expectedOutput: map[string]argusiov1alpha1.ComponentAttestation{},
		},
		{
			name:        "returns attestation",
			cl:          fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeComponentAttestation()).Build(),
			attestation: makeAttestation(),
			expectedOutput: map[string]argusiov1alpha1.ComponentAttestation{
				"test": *makeComponentAttestation(),
			},
		},
		{
			name:          "failed listing",
			cl:            fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			attestation:   makeAttestation(),
			expectedError: "could not list ComponentAttestation",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			output, err := GetComponentAttestations(context.Background(), testCase.cl, testCase.attestation)
			if testCase.expectedError == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedOutput, output)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

func TestLifecycleComponentAttestations(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name          string
		Components    []argusiov1alpha1.ComponentAssessment
		attestations  map[string]argusiov1alpha1.ComponentAttestation
		expectedError string
		AssessmentRef string
		cl            client.Client
	}{
		{
			name:          "no op",
			AssessmentRef: "Assessment",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestations: map[string]argusiov1alpha1.ComponentAttestation{
				"foo": *makeComponentAttestation(),
			},
			Components: []argusiov1alpha1.ComponentAssessment{*makeComponentAssessment()},
		},
		{
			name:          "no Component Assessment label",
			AssessmentRef: "Assessment",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestations: map[string]argusiov1alpha1.ComponentAttestation{
				"foo": *makeComponentAttestation(),
			},
			Components:    []argusiov1alpha1.ComponentAssessment{*makeComponentAssessment(WithLabels(map[string]string{}))},
			expectedError: "Component Assessment 'test' does not contain expected label 'argus.io/Component'",
		},
		{
			name:          "no Component label",
			AssessmentRef: "Assessment",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestations: map[string]argusiov1alpha1.ComponentAttestation{
				"foo": *makeComponentAttestation(AttWithLabels(map[string]string{"argus.io/Assessment": "Assessment"})),
			},
			Components:    []argusiov1alpha1.ComponentAssessment{*makeComponentAssessment()},
			expectedError: "object 'test' does not contain expected label 'argus.io/Component'",
		},
		{
			name:          "Assessment mismatch",
			AssessmentRef: "Assessment2",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeComponentAttestation()).Build(),
			attestations: map[string]argusiov1alpha1.ComponentAttestation{
				"foo": *makeComponentAttestation(AttWithLabels(map[string]string{})),
			},
			Components: []argusiov1alpha1.ComponentAssessment{*makeComponentAssessment()},
		},
		{
			name:          "deletes",
			AssessmentRef: "Assessment",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeComponentAttestation()).Build(),
			attestations: map[string]argusiov1alpha1.ComponentAttestation{
				"foo": *makeComponentAttestation(),
			},
			expectedError: "",
			Components:    []argusiov1alpha1.ComponentAssessment{},
		},
		{
			name:          "fails to delete",
			AssessmentRef: "Assessment",
			cl:            fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			attestations: map[string]argusiov1alpha1.ComponentAttestation{
				"foo": *makeComponentAttestation(),
			},
			Components:    []argusiov1alpha1.ComponentAssessment{},
			expectedError: "could not delete ComponentAttestation",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			err := LifecycleComponentAttestations(context.Background(), testCase.cl, testCase.AssessmentRef, testCase.Components, testCase.attestations)
			if testCase.expectedError == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

func TestCreateOrUpdateComponentAttestations(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name           string
		Components     []argusiov1alpha1.ComponentAssessment
		attestation    *argusiov1alpha1.Attestation
		expectedError  string
		expectedOutput []argusiov1alpha1.NamespacedName
		cl             client.Client
	}{
		{
			name:           "no op",
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestation:    makeAttestation(),
			Components:     []argusiov1alpha1.ComponentAssessment{},
			expectedOutput: []argusiov1alpha1.NamespacedName{},
		},
		{
			name:        "create one",
			cl:          fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestation: makeAttestation(),
			Components:  []argusiov1alpha1.ComponentAssessment{*makeComponentAssessment()},
			expectedOutput: []argusiov1alpha1.NamespacedName{
				{
					Name:      "test-Component",
					Namespace: "default",
				},
			},
		},
		{
			name:          "fail create",
			cl:            fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			attestation:   makeAttestation(),
			Components:    []argusiov1alpha1.ComponentAssessment{*makeComponentAssessment()},
			expectedError: "could not create ComponentAttestation 'test-Component'",
		},
		{
			name:          "fail no Assessment labels",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestation:   makeAttestation(),
			Components:    []argusiov1alpha1.ComponentAssessment{*makeComponentAssessment(WithLabels(map[string]string{}))},
			expectedError: "Component Assessment 'test' does not contain expected label 'argus.io/Assessment'",
		},
		{
			name:          "fail no Component labels",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestation:   makeAttestation(),
			Components:    []argusiov1alpha1.ComponentAssessment{*makeComponentAssessment(WithLabels(map[string]string{"argus.io/Assessment": "Assessment"}))},
			expectedError: "Component Assessment 'test' does not contain expected label 'argus.io/Component'",
		},
		{
			name:          "fail no Control labels",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestation:   makeAttestation(),
			Components:    []argusiov1alpha1.ComponentAssessment{*makeComponentAssessment(WithLabels(map[string]string{"argus.io/Assessment": "Assessment", "argus.io/Component": "Component"}))},
			expectedError: "Component Assessment 'test' does not contain expected label 'argus.io/Control'",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			ouput, err := CreateOrUpdateComponentAttestations(context.Background(), testCase.cl, commonScheme, testCase.attestation, testCase.Components)
			if testCase.expectedError == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedOutput, ouput)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

type attestationMutationFn func(*argusiov1alpha1.Attestation)

func makeAttestation(f ...attestationMutationFn) *argusiov1alpha1.Attestation {
	res := &argusiov1alpha1.Attestation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Attestation",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.AttestationSpec{
			AssessmentRef: "Assessment",
			ProviderRef: argusiov1alpha1.AttestationProviderRef{
				Name:      "bar",
				Namespace: "bing",
			},
		},
	}
	for _, fn := range f {
		fn(res)
	}

	return res
}

type ComponentAttestationMutationFn func(*argusiov1alpha1.ComponentAttestation)

func AttWithLabels(labels map[string]string) ComponentAttestationMutationFn {
	return func(r *argusiov1alpha1.ComponentAttestation) {
		r.Labels = labels
	}
}

func WithName(name string) ComponentAttestationMutationFn {
	return func(r *argusiov1alpha1.ComponentAttestation) {
		r.Name = name
	}
}
func WithResult(result argusiov1alpha1.AttestationResultType) ComponentAttestationMutationFn {
	return func(r *argusiov1alpha1.ComponentAttestation) {
		r.Status.Result.Result = result
	}
}
func makeComponentAttestation(f ...ComponentAttestationMutationFn) *argusiov1alpha1.ComponentAttestation {
	res := &argusiov1alpha1.ComponentAttestation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
			Labels: map[string]string{
				"argus.io/Component":   "Component",
				"argus.io/Assessment":  "Assessment",
				"argus.io/attestation": "test",
			},
			ResourceVersion: "999",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ComponentAttestation",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ComponentAttestationSpec{
			ProviderRef: argusiov1alpha1.AttestationProviderRef{
				Name:      "test",
				Namespace: "test",
			},
		},
		Status: argusiov1alpha1.ComponentAttestationStatus{
			Result: argusiov1alpha1.AttestationResult{
				Result: argusiov1alpha1.AttestationResultTypePass,
			},
		},
	}
	for _, fn := range f {
		fn(res)
	}
	return res
}

type ComponentAssessmentMutationFn func(*argusiov1alpha1.ComponentAssessment)

func WithLabels(labels map[string]string) ComponentAssessmentMutationFn {
	return func(r *argusiov1alpha1.ComponentAssessment) {
		r.Labels = labels
	}
}
func makeComponentAssessment(f ...ComponentAssessmentMutationFn) *argusiov1alpha1.ComponentAssessment {
	res := &argusiov1alpha1.ComponentAssessment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
			Labels: map[string]string{
				"argus.io/Component":  "Component",
				"argus.io/Assessment": "Assessment",
				"argus.io/Control":    "foo_v1",
			},
			ResourceVersion: "999",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ComponentAssessment",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ComponentAssessmentSpec{
			Class: "test",
			ControlRef: argusiov1alpha1.AssessmentControlDefinition{
				Code:    "foo",
				Version: "v1",
			},
		},
	}
	for _, fn := range f {
		fn(res)
	}
	return res
}
