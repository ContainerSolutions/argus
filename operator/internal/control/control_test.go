package control

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

func TestGetComponentControlsFromControl(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name              string
		Control           *argusiov1alpha1.Control
		ComponentControls []argusiov1alpha1.ComponentControl
		expectedOutput    map[string]argusiov1alpha1.ComponentControl
		expectedError     string
		cl                client.Client
	}{
		{
			name: "no Component Controls",
			cl:   fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			Control: &argusiov1alpha1.Control{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "no-Component",
					Namespace: "test",
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "Controls",
					APIVersion: "argus.io/v1alpha1",
				},
				Spec: argusiov1alpha1.ControlSpec{
					Definition: argusiov1alpha1.ControlDefinition{
						Code:    "foo",
						Version: "v1",
					},
				},
			},
			expectedOutput: map[string]argusiov1alpha1.ComponentControl{},
			expectedError:  "",
		},
		{
			name: "Failure Listing",
			cl:   fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			Control: &argusiov1alpha1.Control{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "failure-listing",
					Namespace: "test",
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "Controls",
					APIVersion: "argus.io/v1alpha1",
				},
			},
			expectedOutput: map[string]argusiov1alpha1.ComponentControl{},
			expectedError:  "could not list ComponentControls",
		},
		{
			name: "List And Append",
			cl: fake.NewClientBuilder().WithScheme(commonScheme).WithRuntimeObjects(&argusiov1alpha1.ComponentControl{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ComponentControl",
					APIVersion: "argus.io/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
					Labels: map[string]string{
						"argus.io/Control": "foo_v1",
					},
				},
			}).Build(),
			Control: &argusiov1alpha1.Control{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "list-and-append",
					Namespace: "bar",
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "Controls",
					APIVersion: "argus.io/v1alpha1",
				},
				Spec: argusiov1alpha1.ControlSpec{
					Definition: argusiov1alpha1.ControlDefinition{
						Code:    "foo",
						Version: "v1",
					},
				},
			},
			expectedOutput: map[string]argusiov1alpha1.ComponentControl{
				"foo": {
					ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "bar", ResourceVersion: "999", Labels: map[string]string{"argus.io/Control": "foo_v1"}},
					TypeMeta:   metav1.TypeMeta{Kind: "ComponentControl", APIVersion: "argus.io/v1alpha1"}},
			},
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			output, err := GetComponentControlsFromControl(context.Background(), testCase.cl, testCase.Control)
			if testCase.expectedError == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedOutput, output)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

func TestLifecycleComponentControls(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name              string
		resReq            map[string]argusiov1alpha1.ComponentControl
		Components        []argusiov1alpha1.Component
		classes           []string
		ComponentControls []argusiov1alpha1.ComponentControl
		expectedError     string
		cl                client.Client
	}{
		{
			name:    "no Component Controls",
			cl:      fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			resReq:  map[string]argusiov1alpha1.ComponentControl{},
			classes: []string{"class1"},
			Components: []argusiov1alpha1.Component{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "Component",
						Namespace: "default",
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "Component",
						APIVersion: "argus.io/v1alpha1",
					},
					Spec: argusiov1alpha1.ComponentSpec{
						Parents: []string{"parent"},
						Classes: []string{"class1"},
					},
				},
			},
			expectedError: "",
		},
		{
			name: "Matching Component Controls",
			cl:   fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			resReq: map[string]argusiov1alpha1.ComponentControl{
				"foo": *makeComponentControl(),
			},
			classes:       []string{"class1"},
			Components:    []argusiov1alpha1.Component{*makeComponent()},
			expectedError: "",
		},
		{
			name: "Component Control with no Component",
			cl:   fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeComponentControl()).Build(),
			resReq: map[string]argusiov1alpha1.ComponentControl{
				"foo": *makeComponentControl(),
			},
			classes:       []string{"class1"},
			Components:    []argusiov1alpha1.Component{},
			expectedError: "",
		},
		{
			name: "Component Control with no Component error on delete",
			cl:   fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			resReq: map[string]argusiov1alpha1.ComponentControl{
				"foo": *makeComponentControl(),
			},
			classes:       []string{"class1"},
			Components:    []argusiov1alpha1.Component{},
			expectedError: "could not delete ComponentControl 'foo'",
		},
		{
			name: "No Component Tag",
			cl:   fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeComponentControl()).Build(),
			resReq: map[string]argusiov1alpha1.ComponentControl{
				"foo": *makeComponentControl(WithLabels(map[string]string{})),
			},
			classes:       []string{"class1"},
			Components:    []argusiov1alpha1.Component{*makeComponent()},
			expectedError: "object 'foo' does not contain expected label 'argus.io/Component'",
		},
		{
			name: "No Component-class Tag",
			cl:   fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeComponentControl()).Build(),
			resReq: map[string]argusiov1alpha1.ComponentControl{
				"foo": *makeComponentControl(WithLabels(map[string]string{"argus.io/Component": "foo"})),
			},
			classes:       []string{"class1"},
			Components:    []argusiov1alpha1.Component{*makeComponent()},
			expectedError: "object 'foo' does not contain expected label 'argus.io/Component-class'",
		},
		{
			name: "Control Class changed",
			cl:   fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeComponentControl()).Build(),
			resReq: map[string]argusiov1alpha1.ComponentControl{
				"foo": *makeComponentControl(),
			},
			classes:    []string{"class2"},
			Components: []argusiov1alpha1.Component{*makeComponent()},
		},
		{
			name: "Control Class changed error deleting",
			cl:   fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			resReq: map[string]argusiov1alpha1.ComponentControl{
				"foo": *makeComponentControl(),
			},
			classes:       []string{"class2"},
			Components:    []argusiov1alpha1.Component{*makeComponent()},
			expectedError: "could not delete ComponentControl 'foo'",
		},
		{
			name: "Component Class changed error Deleting",
			cl:   fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			resReq: map[string]argusiov1alpha1.ComponentControl{
				"foo": *makeComponentControl(),
			},
			classes:       []string{"class1"},
			Components:    []argusiov1alpha1.Component{*makeComponent(WithClasses([]string{"class2"}))},
			expectedError: "could not delete ComponentControl 'foo'",
		},
		{
			name: "Component Class changed",
			cl:   fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeComponentControl()).Build(),
			resReq: map[string]argusiov1alpha1.ComponentControl{
				"foo": *makeComponentControl(),
			},
			classes:    []string{"class1"},
			Components: []argusiov1alpha1.Component{*makeComponent(WithClasses([]string{"class2"}))},
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			err := LifecycleComponentControls(context.Background(), testCase.cl, testCase.classes, testCase.Components, testCase.resReq)
			if testCase.expectedError == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

func TestCreateOrUpdateComponentControls(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name          string
		Control       *argusiov1alpha1.Control
		Components    []argusiov1alpha1.Component
		expectedError string
		cl            client.Client
	}{
		{
			name:       "no Controls, no Components",
			cl:         fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			Control:    &argusiov1alpha1.Control{},
			Components: []argusiov1alpha1.Component{},
		},
		{
			name:       "no Components",
			cl:         fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			Control:    makeControl(),
			Components: []argusiov1alpha1.Component{},
		},
		{
			name:       "Component does not contain class",
			cl:         fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			Control:    makeControl(),
			Components: []argusiov1alpha1.Component{*makeComponent(WithClasses([]string{"class2"}))},
		},
		{
			name:       "ComponentControl is created",
			cl:         fake.NewClientBuilder().WithScheme(commonScheme).Build(),
			Control:    makeControl(),
			Components: []argusiov1alpha1.Component{*makeComponent()},
		},
		{
			name:          "ComponentControl creation error",
			cl:            fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			Control:       makeControl(),
			Components:    []argusiov1alpha1.Component{*makeComponent()},
			expectedError: "could not create ComponentControl 'foo-Component'",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			_, err := CreateOrUpdateComponentControls(context.Background(), testCase.cl, commonScheme, testCase.Control, testCase.Components)
			if testCase.expectedError == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

// Helpers

type mutateFunc func(*argusiov1alpha1.ComponentControl)

func WithLabels(labels map[string]string) mutateFunc {
	return func(a *argusiov1alpha1.ComponentControl) {
		a.ObjectMeta.Labels = labels
	}
}
func makeComponentControl(f ...mutateFunc) *argusiov1alpha1.ComponentControl {
	a := &argusiov1alpha1.ComponentControl{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "default",
			Labels: map[string]string{
				"argus.io/Component":       "Component",
				"argus.io/Component-class": "class1",
			},
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ComponentControl",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ComponentControlSpec{
			Definition: argusiov1alpha1.ControlDefinition{
				Code:    "foo",
				Version: "bar",
			},
			RequiredAssessmentClasses: []string{"Assessment"},
		},
	}
	for _, fn := range f {
		fn(a)
	}
	return a
}

type ComponentMutateFn func(*argusiov1alpha1.Component)

func WithClasses(classes []string) ComponentMutateFn {
	return func(a *argusiov1alpha1.Component) {
		a.Spec.Classes = classes
	}
}
func makeComponent(f ...ComponentMutateFn) *argusiov1alpha1.Component {
	a := &argusiov1alpha1.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "Component",
			Namespace: "default",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Component",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ComponentSpec{
			Parents: []string{"parent"},
			Classes: []string{"class1"},
		},
	}
	for _, fn := range f {
		fn(a)
	}
	return a
}

type ControlMutateFn func(*argusiov1alpha1.Control)

func makeControl(f ...ControlMutateFn) *argusiov1alpha1.Control {
	a := &argusiov1alpha1.Control{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "bar",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Control",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ControlSpec{
			Definition: argusiov1alpha1.ControlDefinition{
				Code:    "foo",
				Version: "v1",
			},
			ApplicableComponentClasses: []string{"class1"},
			RequiredAssessmentClasses:  []string{"Assessment"},
		},
	}
	for _, fn := range f {
		fn(a)
	}
	return a
}
