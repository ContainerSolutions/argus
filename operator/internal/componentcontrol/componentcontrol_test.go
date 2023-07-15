package componentcontrol

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

func TestGetValidComponentAssessments(t *testing.T) {
	metrics.SetUpMetrics()
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name          string
		res           *argusiov1alpha1.ComponentControl
		expectedList  []argusiov1alpha1.NamespacedName
		expectedValid int
		expectedError string
		cl            client.Client
	}{
		{
			name:          "No Assessments",
			res:           makeComponentControl(),
			expectedList:  []argusiov1alpha1.NamespacedName{},
			expectedValid: 0,
			expectedError: "",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).Build(),
		},
		{
			name:          "No tags",
			res:           makeComponentControl(WithLabels(map[string]string{})),
			expectedList:  []argusiov1alpha1.NamespacedName{},
			expectedValid: 0,
			expectedError: "object does not have expected label 'argus.io/Component'",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).Build(),
		},
		{
			name: "valid Assessments",
			res:  makeComponentControl(),
			expectedList: []argusiov1alpha1.NamespacedName{
				{
					Name:      "Assessment",
					Namespace: "test",
				},
			},
			expectedValid: 1,
			expectedError: "",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeNewComponentAssessment()).Build(),
		},
		{
			name: "invalid Assessments",
			res:  makeComponentControl(),
			expectedList: []argusiov1alpha1.NamespacedName{
				{
					Name:      "Assessment",
					Namespace: "test",
				},
			},
			expectedValid: 0,
			expectedError: "",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeNewComponentAssessment(WithPass(0))).Build(),
		},
		{
			name:          "Error listing",
			res:           makeComponentControl(),
			expectedList:  []argusiov1alpha1.NamespacedName{},
			expectedValid: 0,
			expectedError: "could not list ComponentAssessment",
			cl:            fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			list, valid, err := GetValidComponentAssessments(context.Background(), testCase.cl, *testCase.res)
			if testCase.expectedError == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedList, list)
				assert.Equal(t, testCase.expectedValid, valid)
			} else {
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

type mutateFn func(*argusiov1alpha1.ComponentControl)

func WithLabels(labels map[string]string) mutateFn {
	return func(res *argusiov1alpha1.ComponentControl) {
		res.Labels = labels
	}
}

func makeComponentControl(f ...mutateFn) *argusiov1alpha1.ComponentControl {
	res := &argusiov1alpha1.ComponentControl{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
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
				Code:    "test",
				Version: "v1",
			},
			RequiredAssessmentClasses: []string{"Assessment"},
		},
	}
	for _, fn := range f {
		fn(res)
	}
	return res
}

type AssessmentFn func(*argusiov1alpha1.ComponentAssessment)

func WithPass(pass int) AssessmentFn {
	return func(res *argusiov1alpha1.ComponentAssessment) {
		res.Status.PassedAttestations = pass
	}
}
func makeNewComponentAssessment(f ...AssessmentFn) *argusiov1alpha1.ComponentAssessment {
	res := &argusiov1alpha1.ComponentAssessment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "Assessment",
			Namespace: "test",
			Labels: map[string]string{
				"argus.io/Component": "Component",
			},
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ComponentAssessment",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ComponentAssessmentSpec{
			Class: "Assessment",
			ControlRef: argusiov1alpha1.AssessmentControlDefinition{
				Code:    "test",
				Version: "v1",
			},
		},
		Status: argusiov1alpha1.ComponentAssessmentStatus{
			TotalAttestations:  1,
			PassedAttestations: 1,
		},
	}
	for _, fn := range f {
		fn(res)
	}
	return res
}
