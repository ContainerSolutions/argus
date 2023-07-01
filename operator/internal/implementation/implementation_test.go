package implementation

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

func TestGetResourceImplementations(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name           string
		implementation *argusiov1alpha1.Implementation
		expectedOutput map[string]argusiov1alpha1.ResourceImplementation
		expectedError  string
		cl             client.Client
	}{
		{
			name:           "fail listing",
			cl:             fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			implementation: makeImplementation(),
			expectedOutput: map[string]argusiov1alpha1.ResourceImplementation{},
			expectedError:  "could not list ResourceImplementation",
		},
		{
			name:           "no Resource Implementations",
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			implementation: makeImplementation(),
			expectedOutput: map[string]argusiov1alpha1.ResourceImplementation{},
			expectedError:  "",
		},
		{
			name:           "with Resource Implementations",
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeResourceImplementation()).Build(),
			implementation: makeImplementation(),
			expectedOutput: map[string]argusiov1alpha1.ResourceImplementation{
				"test": *makeResourceImplementation(),
			},
			expectedError: "",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			output, err := GetResourceImplementations(context.Background(), testCase.cl, testCase.implementation)
			if testCase.expectedError == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedOutput, output)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

func TestLifecycleResourceImplementations(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name          string
		resources     []argusiov1alpha1.Resource
		old           map[string]argusiov1alpha1.ResourceImplementation
		new           map[string]argusiov1alpha1.ResourceImplementation
		cl            client.Client
		expectedError string
	}{
		{
			name:      "noop",
			resources: []argusiov1alpha1.Resource{},
			old:       map[string]argusiov1alpha1.ResourceImplementation{},
			new:       map[string]argusiov1alpha1.ResourceImplementation{},
			cl:        fake.NewClientBuilder().WithScheme(commonScheme).Build(),
		},
		{
			name: "noop with resource",
			old:  map[string]argusiov1alpha1.ResourceImplementation{},
			new: map[string]argusiov1alpha1.ResourceImplementation{
				"test": *makeResourceImplementation()},
			cl: fake.NewClientBuilder().WithScheme(commonScheme).Build(),
		},
		{
			name:      "no resources",
			resources: []argusiov1alpha1.Resource{},
			new:       map[string]argusiov1alpha1.ResourceImplementation{},
			old: map[string]argusiov1alpha1.ResourceImplementation{
				"test": *makeResourceImplementation()},
			cl: fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeResourceImplementation()).Build(),
		},
		{
			name:      "delete failure",
			resources: []argusiov1alpha1.Resource{},
			new:       map[string]argusiov1alpha1.ResourceImplementation{},
			old: map[string]argusiov1alpha1.ResourceImplementation{
				"test": *makeResourceImplementation()},
			cl:            fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			expectedError: "could not delete ResourceImplementation",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			err := LifecycleResourceImplementations(context.Background(), testCase.cl, testCase.new, testCase.old)
			if testCase.expectedError == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}

}
func TestCreateOrUpdateResourceImplementations(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name           string
		implementation *argusiov1alpha1.Implementation
		resources      []argusiov1alpha1.Resource
		expectedError  string
		expectedOutput []argusiov1alpha1.NamespacedName
		cl             client.Client
	}{
		{
			name:           "create one",
			implementation: makeImplementation(),
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			resources:      []argusiov1alpha1.Resource{*makeResource()},
			expectedOutput: []argusiov1alpha1.NamespacedName{
				{
					Name:      "test-resource",
					Namespace: "test",
				},
			},
		},
		{
			name:           "create cascading",
			implementation: makeImplementation(WithCascadePolicy(argusiov1alpha1.CascadingPolicyCascade)),
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			resources:      []argusiov1alpha1.Resource{*makeResource(), *makeResource(WithName("child"))},
			expectedOutput: []argusiov1alpha1.NamespacedName{
				{
					Name:      "test-resource",
					Namespace: "test",
				},
				{
					Name:      "test-child",
					Namespace: "test",
				},
			},
		},
		{
			name:           "fail creating",
			implementation: makeImplementation(),
			cl:             fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			resources:      []argusiov1alpha1.Resource{*makeResource()},
			expectedError:  "could not create resourceImplementation 'test-resource'",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			output, err := CreateOrUpdateResourceImplementations(context.Background(), testCase.cl, commonScheme, testCase.implementation, testCase.resources)
			if testCase.expectedError == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedOutput, output)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

// Helpers

type resourceImplementationMutationFn func(*argusiov1alpha1.ResourceImplementation)

func WithLabels(labels map[string]string) resourceImplementationMutationFn {
	return func(res *argusiov1alpha1.ResourceImplementation) {
		res.ObjectMeta.Labels = labels
	}
}

func makeResourceImplementation(f ...resourceImplementationMutationFn) *argusiov1alpha1.ResourceImplementation {
	res := &argusiov1alpha1.ResourceImplementation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
			Labels: map[string]string{
				"argus.io/resource":       "resource",
				"argus.io/requirement":    "foo_v1",
				"argus.io/implementation": "test",
			},
			ResourceVersion: "999",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ResourceImplementation",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ResourceImplementationSpec{
			Class: "implementation",
			RequirementRef: argusiov1alpha1.ImplementationRequirementDefinition{
				Code:    "foo",
				Version: "v1",
			},
		},
		Status: argusiov1alpha1.ResourceImplementationStatus{
			TotalAttestations:  2,
			PassedAttestations: 2,
		},
	}
	for _, fn := range f {
		fn(res)
	}
	return res
}

type implementationMutationFn func(*argusiov1alpha1.Implementation)

func WithCascadePolicy(policy argusiov1alpha1.ImplementationCascadePolicy) implementationMutationFn {
	return func(res *argusiov1alpha1.Implementation) {
		res.Spec.CascadePolicy = policy
	}
}
func makeImplementation(f ...implementationMutationFn) *argusiov1alpha1.Implementation {
	res := &argusiov1alpha1.Implementation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Implementation",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ImplementationSpec{
			Class: "implementation",
			RequirementRef: argusiov1alpha1.ImplementationRequirementDefinition{
				Code:    "foo",
				Version: "v1",
			},
			ResourceRef: []argusiov1alpha1.NamespacedName{
				{
					Name:      "resource",
					Namespace: "test",
				},
			},
		},
	}
	for _, fn := range f {
		fn(res)
	}
	return res
}

type resourceMutationFn func(*argusiov1alpha1.Resource)

func WithName(name string) resourceMutationFn {
	return func(res *argusiov1alpha1.Resource) {
		res.ObjectMeta.Name = name
	}
}
func makeResource(f ...resourceMutationFn) *argusiov1alpha1.Resource {
	res := &argusiov1alpha1.Resource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "resource",
			Namespace: "test",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Resource",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ResourceSpec{
			Classes: []string{"one", "two"},
			Parents: []string{"parent"},
		},
		Status: argusiov1alpha1.ResourceStatus{
			Children: map[string]argusiov1alpha1.ResourceChild{
				"child": {
					Compliant: true,
				},
			},
		},
	}
	for _, fn := range f {
		fn(res)
	}
	return res
}
