package resource

import (
	"context"
	"testing"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestUpdateRequirements(t *testing.T) {
	testCases := []struct {
		name                         string
		expectedOutput               *argusiov1alpha1.Resource
		inputResource                *argusiov1alpha1.Resource
		inputResourceRequirementList argusiov1alpha1.ResourceRequirementList
	}{
		{
			name:          "No Requirements",
			inputResource: &argusiov1alpha1.Resource{},
			expectedOutput: &argusiov1alpha1.Resource{
				Status: argusiov1alpha1.ResourceStatus{
					Requirements:            map[string]*argusiov1alpha1.ResourceRequirementCompliance{},
					TotalRequirements:       0,
					ImplementedRequirements: 0,
				},
			},
			inputResourceRequirementList: argusiov1alpha1.ResourceRequirementList{
				Items: []argusiov1alpha1.ResourceRequirement{},
			},
		},
		{
			name:          "2 Resource Requirements, no errors",
			inputResource: &argusiov1alpha1.Resource{},
			expectedOutput: &argusiov1alpha1.Resource{
				Status: argusiov1alpha1.ResourceStatus{
					Requirements: map[string]*argusiov1alpha1.ResourceRequirementCompliance{
						"test:1": {
							Implemented: true,
						},
						"test2:1": {
							Implemented: false,
						},
					},
					TotalRequirements:       2,
					ImplementedRequirements: 1,
				},
			},
			inputResourceRequirementList: argusiov1alpha1.ResourceRequirementList{
				Items: []argusiov1alpha1.ResourceRequirement{
					{
						Spec: argusiov1alpha1.ResourceRequirementSpec{
							Definition: argusiov1alpha1.RequirementDefinition{
								Code:    "test",
								Version: "1",
							},
						},
						Status: argusiov1alpha1.ResourceRequirementStatus{
							ValidImplementations: 1,
							TotalImplementations: 1,
						},
					},
					{
						Spec: argusiov1alpha1.ResourceRequirementSpec{
							Definition: argusiov1alpha1.RequirementDefinition{
								Code:    "test2",
								Version: "1",
							},
						},
						Status: argusiov1alpha1.ResourceRequirementStatus{
							ValidImplementations: 1,
							TotalImplementations: 2,
						},
					},
				},
			},
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			output := UpdateRequirements(testCase.inputResourceRequirementList, testCase.inputResource)
			assert.Equal(t, output.Status.Requirements, testCase.expectedOutput.Status.Requirements)
			assert.Equal(t, output.Status.TotalRequirements, testCase.expectedOutput.Status.TotalRequirements)
			assert.Equal(t, output.Status.ImplementedRequirements, testCase.expectedOutput.Status.ImplementedRequirements)
		})
	}
}

func TestUpdateChild(t *testing.T) {
	argusiov1alpha1.AddToScheme(scheme.Scheme)
	testCases := []struct {
		name           string
		expectedOutput string
		inputResource  *argusiov1alpha1.Resource
		clientContent  []client.Object
	}{
		{
			name: "No Parents",
			inputResource: &argusiov1alpha1.Resource{
				Spec: argusiov1alpha1.ResourceSpec{
					Parents: []string{},
				},
			},
		},
		{
			name: "Parent Available",
			inputResource: &argusiov1alpha1.Resource{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: argusiov1alpha1.ResourceSpec{
					Parents: []string{"existing"},
				},
			},
			clientContent: []client.Object{
				&argusiov1alpha1.Resource{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "existing",
						Namespace: "default",
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "Resource",
						APIVersion: "argus.io/v1alpha1",
					},
				},
			},
		},
		{
			name: "Parent Not found",
			inputResource: &argusiov1alpha1.Resource{
				Spec: argusiov1alpha1.ResourceSpec{
					Parents: []string{"unexisting"},
				},
			},
			expectedOutput: "parent resource unexisting not found",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			cl := fake.NewClientBuilder().WithObjects(testCase.clientContent...).WithStatusSubresource(testCase.clientContent...).Build()
			output := UpdateChild(context.TODO(), cl, testCase.inputResource)
			if testCase.expectedOutput == "" {
				assert.Nil(t, output)
			} else {
				assert.ErrorContains(t, output, testCase.expectedOutput)
			}
		})
	}
}
