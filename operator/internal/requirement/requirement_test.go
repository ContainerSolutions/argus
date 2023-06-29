package requirement

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

func TestGetResourceRequirementsFromRequirement(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name                 string
		requirement          *argusiov1alpha1.Requirement
		resourceRequirements []argusiov1alpha1.ResourceRequirement
		expectedOutput       map[string]argusiov1alpha1.ResourceRequirement
		expectedError        string
		cl                   client.Client
	}{
		{
			name: "no Resource Requirements",
			cl:   fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			requirement: &argusiov1alpha1.Requirement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "no-resource",
					Namespace: "test",
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "Requirements",
					APIVersion: "argus.io/v1alpha1",
				},
				Spec: argusiov1alpha1.RequirementSpec{
					Definition: argusiov1alpha1.RequirementDefinition{
						Code:    "foo",
						Version: "v1",
					},
				},
			},
			expectedOutput: map[string]argusiov1alpha1.ResourceRequirement{},
			expectedError:  "",
		},
		{
			name: "Failure Listing",
			cl:   fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			requirement: &argusiov1alpha1.Requirement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "failure-listing",
					Namespace: "test",
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "Requirements",
					APIVersion: "argus.io/v1alpha1",
				},
			},
			expectedOutput: map[string]argusiov1alpha1.ResourceRequirement{},
			expectedError:  "could not list ResourceRequirements",
		},
		{
			name: "List And Append",
			cl: fake.NewClientBuilder().WithScheme(commonScheme).WithRuntimeObjects(&argusiov1alpha1.ResourceRequirement{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ResourceRequirement",
					APIVersion: "argus.io/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
					Labels: map[string]string{
						"argus.io/requirement": "foo_v1",
					},
				},
			}).Build(),
			requirement: &argusiov1alpha1.Requirement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "list-and-append",
					Namespace: "bar",
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "Requirements",
					APIVersion: "argus.io/v1alpha1",
				},
				Spec: argusiov1alpha1.RequirementSpec{
					Definition: argusiov1alpha1.RequirementDefinition{
						Code:    "foo",
						Version: "v1",
					},
				},
			},
			expectedOutput: map[string]argusiov1alpha1.ResourceRequirement{
				"foo": {
					ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "bar", ResourceVersion: "999", Labels: map[string]string{"argus.io/requirement": "foo_v1"}},
					TypeMeta:   metav1.TypeMeta{Kind: "ResourceRequirement", APIVersion: "argus.io/v1alpha1"}},
			},
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			output, err := GetResourceRequirementsFromRequirement(context.Background(), testCase.cl, testCase.requirement)
			if testCase.expectedError == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedOutput, output)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

func TestLifecycleResourceRequirements(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name                 string
		resReq               map[string]argusiov1alpha1.ResourceRequirement
		resources            []argusiov1alpha1.Resource
		classes              []string
		resourceRequirements []argusiov1alpha1.ResourceRequirement
		expectedError        string
		cl                   client.Client
	}{
		{
			name:    "no Resource Requirements",
			cl:      fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			resReq:  map[string]argusiov1alpha1.ResourceRequirement{},
			classes: []string{"class1"},
			resources: []argusiov1alpha1.Resource{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "resource",
						Namespace: "default",
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "Resource",
						APIVersion: "argus.io/v1alpha1",
					},
					Spec: argusiov1alpha1.ResourceSpec{
						Parents: []string{"parent"},
						Classes: []string{"class1"},
					},
				},
			},
			expectedError: "",
		},
		{
			name: "Matching Resource Requirements",
			cl:   fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			resReq: map[string]argusiov1alpha1.ResourceRequirement{
				"foo": *makeResourceRequirement(),
			},
			classes:       []string{"class1"},
			resources:     []argusiov1alpha1.Resource{*makeResource()},
			expectedError: "",
		},
		{
			name: "Resource Requirement with no Resource",
			cl:   fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeResourceRequirement()).Build(),
			resReq: map[string]argusiov1alpha1.ResourceRequirement{
				"foo": *makeResourceRequirement(),
			},
			classes:       []string{"class1"},
			resources:     []argusiov1alpha1.Resource{},
			expectedError: "",
		},
		{
			name: "Resource Requirement with no Resource error on delete",
			cl:   fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			resReq: map[string]argusiov1alpha1.ResourceRequirement{
				"foo": *makeResourceRequirement(),
			},
			classes:       []string{"class1"},
			resources:     []argusiov1alpha1.Resource{},
			expectedError: "could not delete ResourceRequirement 'foo'",
		},
		{
			name: "No Resource Tag",
			cl:   fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeResourceRequirement()).Build(),
			resReq: map[string]argusiov1alpha1.ResourceRequirement{
				"foo": *makeResourceRequirement(WithLabels(map[string]string{})),
			},
			classes:       []string{"class1"},
			resources:     []argusiov1alpha1.Resource{*makeResource()},
			expectedError: "object 'foo' does not contain expected label 'argus.io/resource'",
		},
		{
			name: "No resource-class Tag",
			cl:   fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeResourceRequirement()).Build(),
			resReq: map[string]argusiov1alpha1.ResourceRequirement{
				"foo": *makeResourceRequirement(WithLabels(map[string]string{"argus.io/resource": "foo"})),
			},
			classes:       []string{"class1"},
			resources:     []argusiov1alpha1.Resource{*makeResource()},
			expectedError: "object 'foo' does not contain expected label 'argus.io/resource-class'",
		},
		{
			name: "Requirement Class changed",
			cl:   fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeResourceRequirement()).Build(),
			resReq: map[string]argusiov1alpha1.ResourceRequirement{
				"foo": *makeResourceRequirement(),
			},
			classes:   []string{"class2"},
			resources: []argusiov1alpha1.Resource{*makeResource()},
		},
		{
			name: "Requirement Class changed error deleting",
			cl:   fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			resReq: map[string]argusiov1alpha1.ResourceRequirement{
				"foo": *makeResourceRequirement(),
			},
			classes:       []string{"class2"},
			resources:     []argusiov1alpha1.Resource{*makeResource()},
			expectedError: "could not delete ResourceRequirement 'foo'",
		},
		{
			name: "Resource Class changed error Deleting",
			cl:   fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			resReq: map[string]argusiov1alpha1.ResourceRequirement{
				"foo": *makeResourceRequirement(),
			},
			classes:       []string{"class1"},
			resources:     []argusiov1alpha1.Resource{*makeResource(WithClasses([]string{"class2"}))},
			expectedError: "could not delete ResourceRequirement 'foo'",
		},
		{
			name: "Resource Class changed",
			cl:   fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeResourceRequirement()).Build(),
			resReq: map[string]argusiov1alpha1.ResourceRequirement{
				"foo": *makeResourceRequirement(),
			},
			classes:   []string{"class1"},
			resources: []argusiov1alpha1.Resource{*makeResource(WithClasses([]string{"class2"}))},
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			err := LifecycleResourceRequirements(context.Background(), testCase.cl, testCase.classes, testCase.resources, testCase.resReq)
			if testCase.expectedError == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

type mutateFunc func(*argusiov1alpha1.ResourceRequirement)

func WithLabels(labels map[string]string) mutateFunc {
	return func(a *argusiov1alpha1.ResourceRequirement) {
		a.ObjectMeta.Labels = labels
	}
}
func makeResourceRequirement(f ...mutateFunc) *argusiov1alpha1.ResourceRequirement {
	a := &argusiov1alpha1.ResourceRequirement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "default",
			Labels: map[string]string{
				"argus.io/resource":       "resource",
				"argus.io/resource-class": "class1",
			},
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ResourceRequirement",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ResourceRequirementSpec{
			Definition: argusiov1alpha1.RequirementDefinition{
				Code:    "foo",
				Version: "bar",
			},
			RequiredImplementationClasses: []string{"implementation"},
		},
	}
	for _, fn := range f {
		fn(a)
	}
	return a
}

type resourceMutateFn func(*argusiov1alpha1.Resource)

func WithClasses(classes []string) resourceMutateFn {
	return func(a *argusiov1alpha1.Resource) {
		a.Spec.Classes = classes
	}
}
func makeResource(f ...resourceMutateFn) *argusiov1alpha1.Resource {
	a := &argusiov1alpha1.Resource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "resource",
			Namespace: "default",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Resource",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ResourceSpec{
			Parents: []string{"parent"},
			Classes: []string{"class1"},
		},
	}
	for _, fn := range f {
		fn(a)
	}
	return a
}
