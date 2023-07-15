package componentassessment

import (
	"context"
	"testing"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/ContainerSolutions/argus/operator/internal/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestListComponentAttestations(t *testing.T) {
	commonScheme := runtime.NewScheme()
	metrics.SetUpMetrics()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name           string
		Assessment     *argusiov1alpha1.ComponentAssessment
		expectedError  string
		expectedOutput []argusiov1alpha1.ComponentAttestation
		cl             client.Client
	}{
		{
			name:           "no Component attestation",
			Assessment:     makeComponentAssessment(),
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			expectedOutput: []argusiov1alpha1.ComponentAttestation{},
			expectedError:  "",
		},
		{
			name:           "one Component attestation",
			Assessment:     makeComponentAssessment(),
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeComponentAttestation()).Build(),
			expectedOutput: []argusiov1alpha1.ComponentAttestation{*makeComponentAttestation()},
			expectedError:  "",
		},
		{
			name:           "no Component label error",
			Assessment:     makeComponentAssessment(WithLabels(map[string]string{"argus.io/Assessment": "Assessment"})),
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			expectedOutput: []argusiov1alpha1.ComponentAttestation{},
			expectedError:  "object does not have expected label 'argus.io/Component'",
		},
		{
			name:           "no Assessment label error",
			Assessment:     makeComponentAssessment(WithLabels(map[string]string{"argus.io/Component": "Component"})),
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			expectedOutput: []argusiov1alpha1.ComponentAttestation{},
			expectedError:  "object does not have expected label 'argus.io/Assessment'",
		},
		{
			name:          "fail listing",
			Assessment:    makeComponentAssessment(),
			cl:            fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			expectedError: "could not list ComponentAssessment",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			output, err := ListComponentAttestations(context.Background(), testCase.cl, testCase.Assessment)
			if testCase.expectedError == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedOutput, output)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

func TestGetValidComponentAttestations(t *testing.T) {
	testCases := []struct {
		name           string
		attestations   []argusiov1alpha1.ComponentAttestation
		expectedOutput []argusiov1alpha1.NamespacedName
		expectedValid  int
	}{
		{
			name:           "no Component attestation",
			attestations:   []argusiov1alpha1.ComponentAttestation{},
			expectedOutput: []argusiov1alpha1.NamespacedName{},
			expectedValid:  0,
		},
		{
			name:         "one valid Component attestation",
			attestations: []argusiov1alpha1.ComponentAttestation{*makeComponentAttestation()},
			expectedOutput: []argusiov1alpha1.NamespacedName{
				{
					Name:      "test",
					Namespace: "test",
				},
			},
			expectedValid: 1,
		},
		{
			name:         "two manifests, one valid Component attestation",
			attestations: []argusiov1alpha1.ComponentAttestation{*makeComponentAttestation(), *makeComponentAttestation(WithName("test2"), WithResult(argusiov1alpha1.AttestationResultTypeFail))},
			expectedOutput: []argusiov1alpha1.NamespacedName{
				{
					Name:      "test",
					Namespace: "test",
				},
				{
					Name:      "test2",
					Namespace: "test",
				},
			},
			expectedValid: 1,
		},
		{
			name:           "one valid, one invalid Component attestation",
			attestations:   []argusiov1alpha1.ComponentAttestation{},
			expectedOutput: []argusiov1alpha1.NamespacedName{},
			expectedValid:  0,
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			output, valid := GetValidComponentAttestations(context.Background(), testCase.attestations)
			assert.Equal(t, testCase.expectedOutput, output)
			assert.Equal(t, testCase.expectedValid, valid)
		})
	}
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

type ComponentAttestationMutationFn func(*argusiov1alpha1.ComponentAttestation)

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
				"argus.io/Component":  "Component",
				"argus.io/Assessment": "Assessment",
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
