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
