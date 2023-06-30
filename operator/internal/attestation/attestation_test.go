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

func TestGetResourceAttestations(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name           string
		attestation    *argusiov1alpha1.Attestation
		expectedOutput map[string]argusiov1alpha1.ResourceAttestation
		expectedError  string
		cl             client.Client
	}{
		{
			name:           "no op",
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestation:    makeAttestation(),
			expectedOutput: map[string]argusiov1alpha1.ResourceAttestation{},
		},
		{
			name:        "returns attestation",
			cl:          fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeResourceAttestation()).Build(),
			attestation: makeAttestation(),
			expectedOutput: map[string]argusiov1alpha1.ResourceAttestation{
				"test": *makeResourceAttestation(),
			},
		},
		{
			name:          "failed listing",
			cl:            fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			attestation:   makeAttestation(),
			expectedError: "could not list ResourceAttestation",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			output, err := GetResourceAttestations(context.Background(), testCase.cl, testCase.attestation)
			if testCase.expectedError == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedOutput, output)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

func TestLifecycleResourceAttestations(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name              string
		resources         []argusiov1alpha1.ResourceImplementation
		attestations      map[string]argusiov1alpha1.ResourceAttestation
		expectedError     string
		implementationRef string
		cl                client.Client
	}{
		{
			name:              "no op",
			implementationRef: "implementation",
			cl:                fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestations: map[string]argusiov1alpha1.ResourceAttestation{
				"foo": *makeResourceAttestation(),
			},
			resources: []argusiov1alpha1.ResourceImplementation{*makeResourceImplementation()},
		},
		{
			name:              "no resource implementation label",
			implementationRef: "implementation",
			cl:                fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestations: map[string]argusiov1alpha1.ResourceAttestation{
				"foo": *makeResourceAttestation(),
			},
			resources:     []argusiov1alpha1.ResourceImplementation{*makeResourceImplementation(WithLabels(map[string]string{}))},
			expectedError: "resource implementation 'test' does not contain expected label 'argus.io/resource'",
		},
		{
			name:              "no resource label",
			implementationRef: "implementation",
			cl:                fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestations: map[string]argusiov1alpha1.ResourceAttestation{
				"foo": *makeResourceAttestation(AttWithLabels(map[string]string{"argus.io/implementation": "implementation"})),
			},
			resources:     []argusiov1alpha1.ResourceImplementation{*makeResourceImplementation()},
			expectedError: "object 'test' does not contain expected label 'argus.io/resource'",
		},
		{
			name:              "implementation mismatch",
			implementationRef: "implementation2",
			cl:                fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeResourceAttestation()).Build(),
			attestations: map[string]argusiov1alpha1.ResourceAttestation{
				"foo": *makeResourceAttestation(AttWithLabels(map[string]string{})),
			},
			resources: []argusiov1alpha1.ResourceImplementation{*makeResourceImplementation()},
		},
		{
			name:              "deletes",
			implementationRef: "implementation",
			cl:                fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeResourceAttestation()).Build(),
			attestations: map[string]argusiov1alpha1.ResourceAttestation{
				"foo": *makeResourceAttestation(),
			},
			expectedError: "",
			resources:     []argusiov1alpha1.ResourceImplementation{},
		},
		{
			name:              "fails to delete",
			implementationRef: "implementation",
			cl:                fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			attestations: map[string]argusiov1alpha1.ResourceAttestation{
				"foo": *makeResourceAttestation(),
			},
			resources:     []argusiov1alpha1.ResourceImplementation{},
			expectedError: "could not delete ResourceAttestation",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			err := LifecycleResourceAttestations(context.Background(), testCase.cl, testCase.implementationRef, testCase.resources, testCase.attestations)
			if testCase.expectedError == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

func TestCreateOrUpdateResourceAttestations(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name           string
		resources      []argusiov1alpha1.ResourceImplementation
		attestation    *argusiov1alpha1.Attestation
		expectedError  string
		expectedOutput []argusiov1alpha1.NamespacedName
		cl             client.Client
	}{
		{
			name:           "no op",
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestation:    makeAttestation(),
			resources:      []argusiov1alpha1.ResourceImplementation{},
			expectedOutput: []argusiov1alpha1.NamespacedName{},
		},
		{
			name:        "create one",
			cl:          fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestation: makeAttestation(),
			resources:   []argusiov1alpha1.ResourceImplementation{*makeResourceImplementation()},
			expectedOutput: []argusiov1alpha1.NamespacedName{
				{
					Name:      "test-resource",
					Namespace: "default",
				},
			},
		},
		{
			name:          "fail create",
			cl:            fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			attestation:   makeAttestation(),
			resources:     []argusiov1alpha1.ResourceImplementation{*makeResourceImplementation()},
			expectedError: "could not create ResourceAttestation 'test-resource'",
		},
		{
			name:          "fail no implementation labels",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestation:   makeAttestation(),
			resources:     []argusiov1alpha1.ResourceImplementation{*makeResourceImplementation(WithLabels(map[string]string{}))},
			expectedError: "resource implementation 'test' does not contain expected label 'argus.io/implementation'",
		},
		{
			name:          "fail no resource labels",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestation:   makeAttestation(),
			resources:     []argusiov1alpha1.ResourceImplementation{*makeResourceImplementation(WithLabels(map[string]string{"argus.io/implementation": "implementation"}))},
			expectedError: "resource implementation 'test' does not contain expected label 'argus.io/resource'",
		},
		{
			name:          "fail no requirement labels",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			attestation:   makeAttestation(),
			resources:     []argusiov1alpha1.ResourceImplementation{*makeResourceImplementation(WithLabels(map[string]string{"argus.io/implementation": "implementation", "argus.io/resource": "resource"}))},
			expectedError: "resource implementation 'test' does not contain expected label 'argus.io/requirement'",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			ouput, err := CreateOrUpdateResourceAttestations(context.Background(), testCase.cl, testCase.attestation, testCase.resources)
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
			ImplementationRef: "implementation",
			Type: argusiov1alpha1.AttestationType{
				Kind:      "foo",
				Name:      "bar",
				Namespace: "weird",
			},
			ProviderRef: argusiov1alpha1.AttestationProvider{
				Kind:      "foo",
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

type resourceAttestationMutationFn func(*argusiov1alpha1.ResourceAttestation)

func AttWithLabels(labels map[string]string) resourceAttestationMutationFn {
	return func(r *argusiov1alpha1.ResourceAttestation) {
		r.Labels = labels
	}
}

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
				"argus.io/attestation":    "test",
			},
			ResourceVersion: "999",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ResourceAttestation",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ResourceAttestationSpec{
			ProviderRef: argusiov1alpha1.AttestationProvider{
				Name:      "test",
				Kind:      "fake",
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
				"argus.io/requirement":    "foo_v1",
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
