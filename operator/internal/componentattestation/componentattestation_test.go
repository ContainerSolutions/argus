package componentattestation

import (
	"context"
	"fmt"
	"testing"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/ContainerSolutions/argus/operator/internal/provider/schema"
)

type MockClient struct {
	AttestFn func() (argusiov1alpha1.AttestationResult, error)
	CloseFn  func() error
}

func (m *MockClient) Attest() (argusiov1alpha1.AttestationResult, error) {
	return m.AttestFn()
}

func (m *MockClient) Close() error {
	return m.CloseFn()
}

type MockProvider struct {
	NewFn func(*argusiov1alpha1.AttestationProviderSpec) (schema.AttestationClient, error)
}

func (m *MockProvider) New(a *argusiov1alpha1.AttestationProviderSpec) (schema.AttestationClient, error) {
	return m.NewFn(a)
}

func TestGetAttestationClient(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name          string
		prov          MockProvider
		NewFn         func(*argusiov1alpha1.AttestationProviderSpec) (schema.AttestationClient, error)
		spec          *argusiov1alpha1.ComponentAttestation
		cl            client.Client
		expectedError string
	}{
		{
			name:  "success",
			prov:  MockProvider{},
			NewFn: DefaultNewFn(),
			spec:  makeComponentAttestation(),
			cl:    fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeAttestationProvider()).Build(),
		},
		{
			name:          "Provider not found",
			prov:          MockProvider{},
			NewFn:         DefaultNewFn(),
			spec:          makeComponentAttestation(),
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			expectedError: "could not get provider spec 'prov'",
		},
		{
			name:          "Provider not loaded",
			prov:          MockProvider{},
			NewFn:         DefaultNewFn(),
			spec:          makeComponentAttestation(),
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeAttestationProvider(WithType("mack"))).Build(),
			expectedError: "could not get provider 'prov'",
		},
		{
			name:          "Provider New() fails",
			prov:          MockProvider{},
			NewFn:         WithNewFn(nil, fmt.Errorf("boom")),
			spec:          makeComponentAttestation(),
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeAttestationProvider()).Build(),
			expectedError: "could not instantiate client for provider 'prov'",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		schema.ForceRegister(&testCase.prov, "mock")
		t.Run(testCase.name, func(t *testing.T) {
			testCase.prov.NewFn = testCase.NewFn
			_, err := GetAttestationClient(context.Background(), testCase.cl, testCase.spec)
			if testCase.expectedError == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

type ProvFn func(*argusiov1alpha1.AttestationProvider)

func DefaultNewFn() func(*argusiov1alpha1.AttestationProviderSpec) (schema.AttestationClient, error) {
	return func(*argusiov1alpha1.AttestationProviderSpec) (schema.AttestationClient, error) {
		return &MockClient{}, nil
	}
}
func WithNewFn(s schema.AttestationClient, err error) func(*argusiov1alpha1.AttestationProviderSpec) (schema.AttestationClient, error) {
	return func(*argusiov1alpha1.AttestationProviderSpec) (schema.AttestationClient, error) {
		return s, err
	}
}
func WithType(t string) ProvFn {
	return func(p *argusiov1alpha1.AttestationProvider) {
		p.Spec.Type = t
	}
}
func makeAttestationProvider(f ...ProvFn) *argusiov1alpha1.AttestationProvider {
	res := &argusiov1alpha1.AttestationProvider{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prov",
			Namespace: "prov",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "AttestationProvider",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.AttestationProviderSpec{
			Type:           "mock",
			ProviderConfig: map[string]string{},
		},
	}
	for _, m := range f {
		m(res)
	}
	return res
}

type MutationFn func(*argusiov1alpha1.ComponentAttestation)

func makeComponentAttestation(f ...MutationFn) *argusiov1alpha1.ComponentAttestation {
	res := &argusiov1alpha1.ComponentAttestation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ComponentAttestation",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ComponentAttestationSpec{
			ProviderRef: argusiov1alpha1.AttestationProviderRef{
				Name:      "prov",
				Namespace: "prov",
			},
		},
	}

	for _, m := range f {
		m(res)
	}
	return res
}
