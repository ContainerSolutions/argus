package resourcerequirement

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

func TestGetValidResourceImplementations(t *testing.T) {
	commonScheme := runtime.NewScheme()
	err := argusiov1alpha1.AddToScheme(commonScheme)
	require.Nil(t, err)
	testCases := []struct {
		name          string
		res           *argusiov1alpha1.ResourceRequirement
		expectedList  []argusiov1alpha1.NamespacedName
		expectedValid int
		expectedError string
		cl            client.Client
	}{
		{
			name:          "No implementations",
			res:           makeResourceRequirement(),
			expectedList:  []argusiov1alpha1.NamespacedName{},
			expectedValid: 0,
			expectedError: "",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).Build(),
		},
		{
			name:          "No tags",
			res:           makeResourceRequirement(WithLabels(map[string]string{})),
			expectedList:  []argusiov1alpha1.NamespacedName{},
			expectedValid: 0,
			expectedError: "object does not have expected label 'argus.io/resource'",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).Build(),
		},
		{
			name: "valid implementations",
			res:  makeResourceRequirement(),
			expectedList: []argusiov1alpha1.NamespacedName{
				{
					Name:      "implementation",
					Namespace: "test",
				},
			},
			expectedValid: 1,
			expectedError: "",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeNewResourceImplementation()).Build(),
		},
		{
			name: "invalid implementations",
			res:  makeResourceRequirement(),
			expectedList: []argusiov1alpha1.NamespacedName{
				{
					Name:      "implementation",
					Namespace: "test",
				},
			},
			expectedValid: 0,
			expectedError: "",
			cl:            fake.NewClientBuilder().WithScheme(commonScheme).WithObjects(makeNewResourceImplementation(WithPass(0))).Build(),
		},
		{
			name:          "Error listing",
			res:           makeResourceRequirement(),
			expectedList:  []argusiov1alpha1.NamespacedName{},
			expectedValid: 0,
			expectedError: "could not list ResourceImplementation",
			cl:            fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			list, valid, err := GetValidResourceImplementations(context.Background(), testCase.cl, *testCase.res)
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

type mutateFn func(*argusiov1alpha1.ResourceRequirement)

func WithLabels(labels map[string]string) mutateFn {
	return func(res *argusiov1alpha1.ResourceRequirement) {
		res.Labels = labels
	}
}

func makeResourceRequirement(f ...mutateFn) *argusiov1alpha1.ResourceRequirement {
	res := &argusiov1alpha1.ResourceRequirement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
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
				Code:    "test",
				Version: "v1",
			},
			RequiredImplementationClasses: []string{"implementation"},
		},
	}
	for _, fn := range f {
		fn(res)
	}
	return res
}

type implementationFn func(*argusiov1alpha1.ResourceImplementation)

func WithPass(pass int) implementationFn {
	return func(res *argusiov1alpha1.ResourceImplementation) {
		res.Status.PassedAttestations = pass
	}
}
func makeNewResourceImplementation(f ...implementationFn) *argusiov1alpha1.ResourceImplementation {
	res := &argusiov1alpha1.ResourceImplementation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "implementation",
			Namespace: "test",
			Labels: map[string]string{
				"argus.io/resource": "resource",
			},
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ResourceImplementation",
			APIVersion: "argus.io/v1alpha1",
		},
		Spec: argusiov1alpha1.ResourceImplementationSpec{
			RequirementRef: argusiov1alpha1.ImplementationRequirementDefinition{
				Code:    "test",
				Version: "v1",
			},
			ResourceRef: argusiov1alpha1.NamespacedName{
				Name:      "resource",
				Namespace: "test",
			},
		},
		Status: argusiov1alpha1.ResourceImplementationStatus{
			TotalAttestations:  1,
			PassedAttestations: 1,
		},
	}
	for _, fn := range f {
		fn(res)
	}
	return res
}
