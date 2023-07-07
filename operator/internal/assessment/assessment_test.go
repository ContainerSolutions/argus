package assessment

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

func TestGetComponentAssessments(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name           string
		Assessment     *argusiov1alpha1.Assessment
		expectedOutput map[string]argusiov1alpha1.ComponentAssessment
		expectedError  string
		cl             client.Client
	}{
		{
			name:           "fail listing",
			cl:             fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			Assessment:     makeAssessment(),
			expectedOutput: map[string]argusiov1alpha1.ComponentAssessment{},
			expectedError:  "could not list ComponentAssessment",
		},
		{
			name:           "no Component Assessments",
			cl:             fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			Assessment:     makeAssessment(),
			expectedOutput: map[string]argusiov1alpha1.ComponentAssessment{},
			expectedError:  "",
		},
		{
			name:       "with Component Assessments",
			cl:         fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeComponentAssessment()).Build(),
			Assessment: makeAssessment(),
			expectedOutput: map[string]argusiov1alpha1.ComponentAssessment{
				"test": *makeComponentAssessment(),
			},
			expectedError: "",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			output, err := GetComponentAssessments(context.Background(), testCase.cl, testCase.Assessment)
			if testCase.expectedError == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedOutput, output)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

func TestLifecycleComponentAssessments(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name          string
		Components    []argusiov1alpha1.Component
		old           map[string]argusiov1alpha1.ComponentAssessment
		new           map[string]argusiov1alpha1.ComponentAssessment
		cl            client.Client
		expectedError string
	}{
		{
			name:       "noop",
			Components: []argusiov1alpha1.Component{},
			old:        map[string]argusiov1alpha1.ComponentAssessment{},
			new:        map[string]argusiov1alpha1.ComponentAssessment{},
			cl:         fake.NewClientBuilder().WithScheme(commonScheme).Build(),
		},
		{
			name: "noop with Component",
			old:  map[string]argusiov1alpha1.ComponentAssessment{},
			new: map[string]argusiov1alpha1.ComponentAssessment{
				"test": *makeComponentAssessment()},
			cl: fake.NewClientBuilder().WithScheme(commonScheme).Build(),
		},
		{
			name:       "no Components",
			Components: []argusiov1alpha1.Component{},
			new:        map[string]argusiov1alpha1.ComponentAssessment{},
			old: map[string]argusiov1alpha1.ComponentAssessment{
				"test": *makeComponentAssessment()},
			cl: fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeComponentAssessment()).Build(),
		},
		{
			name:       "delete failure",
			Components: []argusiov1alpha1.Component{},
			new:        map[string]argusiov1alpha1.ComponentAssessment{},
			old: map[string]argusiov1alpha1.ComponentAssessment{
				"test": *makeComponentAssessment()},
			cl:            fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			expectedError: "could not delete ComponentAssessment",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			err := LifecycleComponentAssessments(context.Background(), testCase.cl, testCase.new, testCase.old)
			if testCase.expectedError == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}

}
func TestCreateOrUpdateComponentAssessments(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name           string
		Assessment     *argusiov1alpha1.Assessment
		Components     []argusiov1alpha1.Component
		expectedError  string
		expectedOutput []argusiov1alpha1.NamespacedName
		cl             client.Client
	}{
		{
			name:       "create one",
			Assessment: makeAssessment(),
			cl:         fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			Components: []argusiov1alpha1.Component{*makeComponent()},
			expectedOutput: []argusiov1alpha1.NamespacedName{
				{
					Name:      "test-Component",
					Namespace: "test",
				},
			},
		},
		{
			name:       "create cascading",
			Assessment: makeAssessment(WithCascadePolicy(argusiov1alpha1.CascadingPolicyCascade)),
			cl:         fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			Components: []argusiov1alpha1.Component{*makeComponent(), *makeComponent(WithName("child"))},
			expectedOutput: []argusiov1alpha1.NamespacedName{
				{
					Name:      "test-Component",
					Namespace: "test",
				},
				{
					Name:      "test-child",
					Namespace: "test",
				},
			},
		},
		{
			name:          "fail creating",
			Assessment:    makeAssessment(),
			cl:            fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			Components:    []argusiov1alpha1.Component{*makeComponent()},
			expectedError: "could not create ComponentAssessment 'test-Component'",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			output, err := CreateOrUpdateComponentAssessments(context.Background(), testCase.cl, commonScheme, testCase.Assessment, testCase.Components)
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

type ComponentAssessmentMutationFn func(*argusiov1alpha1.ComponentAssessment)

func WithLabels(labels map[string]string) ComponentAssessmentMutationFn {
	return func(res *argusiov1alpha1.ComponentAssessment) {
		res.ObjectMeta.Labels = labels
	}
}

func makeComponentAssessment(f ...ComponentAssessmentMutationFn) *argusiov1alpha1.ComponentAssessment {
	res := &argusiov1alpha1.ComponentAssessment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
			Labels: map[string]string{
				"argus.io/Component":  "Component",
				"argus.io/Control":    "foo_v1",
				"argus.io/Assessment": "test",
			},
			ResourceVersion: "999",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ComponentAssessment",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ComponentAssessmentSpec{
			Class: "Assessment",
			ControlRef: argusiov1alpha1.AssessmentControlDefinition{
				Code:    "foo",
				Version: "v1",
			},
		},
		Status: argusiov1alpha1.ComponentAssessmentStatus{
			TotalAttestations:  2,
			PassedAttestations: 2,
		},
	}
	for _, fn := range f {
		fn(res)
	}
	return res
}

type AssessmentMutationFn func(*argusiov1alpha1.Assessment)

func WithCascadePolicy(policy argusiov1alpha1.AssessmentCascadePolicy) AssessmentMutationFn {
	return func(res *argusiov1alpha1.Assessment) {
		res.Spec.CascadePolicy = policy
	}
}
func makeAssessment(f ...AssessmentMutationFn) *argusiov1alpha1.Assessment {
	res := &argusiov1alpha1.Assessment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Assessment",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.AssessmentSpec{
			Class: "Assessment",
			ControlRef: argusiov1alpha1.AssessmentControlDefinition{
				Code:    "foo",
				Version: "v1",
			},
			ComponentRef: []argusiov1alpha1.NamespacedName{
				{
					Name:      "Component",
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

type ComponentMutationFn func(*argusiov1alpha1.Component)

func WithName(name string) ComponentMutationFn {
	return func(res *argusiov1alpha1.Component) {
		res.ObjectMeta.Name = name
	}
}
func makeComponent(f ...ComponentMutationFn) *argusiov1alpha1.Component {
	res := &argusiov1alpha1.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "Component",
			Namespace: "test",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Component",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ComponentSpec{
			Classes: []string{"one", "two"},
			Parents: []string{"parent"},
		},
		Status: argusiov1alpha1.ComponentStatus{
			Children: map[string]argusiov1alpha1.ComponentChild{
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
