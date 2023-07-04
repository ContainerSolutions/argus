package resourceimplementation

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

func TestListResourceAttestations(t *testing.T) {
	commonScheme := runtime.NewScheme()
	metrics.SetUpMetrics()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name           string
		implementation *argusiov1alpha1.ResourceImplementation
		expectedError  string
		expectedOutput []argusiov1alpha1.ResourceAttestation
		cl             client.Client
	}{
		{
			name:           "no resource attestation",
			implementation: makeResourceImplementation(),
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			expectedOutput: []argusiov1alpha1.ResourceAttestation{},
			expectedError:  "",
		},
		{
			name:           "one resource attestation",
			implementation: makeResourceImplementation(),
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeResourceAttestation()).Build(),
			expectedOutput: []argusiov1alpha1.ResourceAttestation{*makeResourceAttestation()},
			expectedError:  "",
		},
		{
			name:           "no resource label error",
			implementation: makeResourceImplementation(WithLabels(map[string]string{"argus.io/implementation": "implementation"})),
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			expectedOutput: []argusiov1alpha1.ResourceAttestation{},
			expectedError:  "object does not have expected label 'argus.io/resource'",
		},
		{
			name:           "no implementation label error",
			implementation: makeResourceImplementation(WithLabels(map[string]string{"argus.io/resource": "resource"})),
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			expectedOutput: []argusiov1alpha1.ResourceAttestation{},
			expectedError:  "object does not have expected label 'argus.io/implementation'",
		},
		{
			name:           "fail listing",
			implementation: makeResourceImplementation(),
			cl:             fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			expectedError:  "could not list ResourceImplementation",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			output, err := ListResourceAttestations(context.Background(), testCase.cl, testCase.implementation)
			if testCase.expectedError == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedOutput, output)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

func TestGetValidResourceAttestations(t *testing.T) {
	testCases := []struct {
		name           string
		attestations   []argusiov1alpha1.ResourceAttestation
		expectedOutput []argusiov1alpha1.NamespacedName
		expectedValid  int
	}{
		{
			name:           "no resource attestation",
			attestations:   []argusiov1alpha1.ResourceAttestation{},
			expectedOutput: []argusiov1alpha1.NamespacedName{},
			expectedValid:  0,
		},
		{
			name:         "one valid resource attestation",
			attestations: []argusiov1alpha1.ResourceAttestation{*makeResourceAttestation()},
			expectedOutput: []argusiov1alpha1.NamespacedName{
				{
					Name:      "test",
					Namespace: "test",
				},
			},
			expectedValid: 1,
		},
		{
			name:         "two manifests, one valid resource attestation",
			attestations: []argusiov1alpha1.ResourceAttestation{*makeResourceAttestation(), *makeResourceAttestation(WithName("test2"), WithResult(argusiov1alpha1.AttestationResultTypeFail))},
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
			name:           "one valid, one invalid resource attestation",
			attestations:   []argusiov1alpha1.ResourceAttestation{},
			expectedOutput: []argusiov1alpha1.NamespacedName{},
			expectedValid:  0,
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			output, valid := GetValidResourceAttestations(context.Background(), testCase.attestations)
			assert.Equal(t, testCase.expectedOutput, output)
			assert.Equal(t, testCase.expectedValid, valid)
		})
	}
}

type resourceImplementationMutationFn func(*argusiov1alpha1.ResourceImplementation)

func WithLabels(labels map[string]string) resourceImplementationMutationFn {
	return func(r *argusiov1alpha1.ResourceImplementation) {
		r.Labels = labels
	}
}
func makeResourceImplementation(f ...resourceImplementationMutationFn) *argusiov1alpha1.ResourceImplementation {
	res := &argusiov1alpha1.ResourceImplementation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
			Labels: map[string]string{
				"argus.io/resource":       "resource",
				"argus.io/implementation": "implementation",
			},
			ResourceVersion: "999",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ResourceImplementation",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ResourceImplementationSpec{
			Class: "test",
			RequirementRef: argusiov1alpha1.ImplementationRequirementDefinition{
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

type resourceAttestationMutationFn func(*argusiov1alpha1.ResourceAttestation)

func WithName(name string) resourceAttestationMutationFn {
	return func(r *argusiov1alpha1.ResourceAttestation) {
		r.Name = name
	}
}
func WithResult(result argusiov1alpha1.AttestationResultType) resourceAttestationMutationFn {
	return func(r *argusiov1alpha1.ResourceAttestation) {
		r.Status.Result.Result = result
	}
}
func makeResourceAttestation(f ...resourceAttestationMutationFn) *argusiov1alpha1.ResourceAttestation {
	res := &argusiov1alpha1.ResourceAttestation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
			Labels: map[string]string{
				"argus.io/resource":       "resource",
				"argus.io/implementation": "implementation",
			},
			ResourceVersion: "999",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ResourceAttestation",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ResourceAttestationSpec{
			ProviderRef: argusiov1alpha1.AttestationProviderRef{
				Name:      "test",
				Namespace: "test",
			},
		},
		Status: argusiov1alpha1.ResourceAttestationStatus{
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
