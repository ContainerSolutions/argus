package component

import (
	"context"
	"testing"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/ContainerSolutions/argus/operator/internal/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestUpdateControls(t *testing.T) {
	metrics.SetUpMetrics()
	testCases := []struct {
		name                      string
		expectedOutput            *argusiov1alpha1.Component
		inputComponent            *argusiov1alpha1.Component
		inputComponentControlList argusiov1alpha1.ComponentControlList
	}{
		{
			name:           "No Controls",
			inputComponent: &argusiov1alpha1.Component{},
			expectedOutput: &argusiov1alpha1.Component{
				Status: argusiov1alpha1.ComponentStatus{
					Controls:            map[string]*argusiov1alpha1.ComponentControlCompliance{},
					TotalControls:       0,
					ImplementedControls: 0,
				},
			},
			inputComponentControlList: argusiov1alpha1.ComponentControlList{
				Items: []argusiov1alpha1.ComponentControl{},
			},
		},
		{
			name:           "2 Component Controls, no errors",
			inputComponent: &argusiov1alpha1.Component{},
			expectedOutput: &argusiov1alpha1.Component{
				Status: argusiov1alpha1.ComponentStatus{
					Controls: map[string]*argusiov1alpha1.ComponentControlCompliance{
						"test:1": {
							Implemented: true,
						},
						"test2:1": {
							Implemented: false,
						},
					},
					TotalControls:       2,
					ImplementedControls: 1,
				},
			},
			inputComponentControlList: argusiov1alpha1.ComponentControlList{
				Items: []argusiov1alpha1.ComponentControl{
					{
						Spec: argusiov1alpha1.ComponentControlSpec{
							Definition: argusiov1alpha1.ControlDefinition{
								Code:    "test",
								Version: "1",
							},
						},
						Status: argusiov1alpha1.ComponentControlStatus{
							ValidAssessments: 1,
							TotalAssessments: 1,
						},
					},
					{
						Spec: argusiov1alpha1.ComponentControlSpec{
							Definition: argusiov1alpha1.ControlDefinition{
								Code:    "test2",
								Version: "1",
							},
						},
						Status: argusiov1alpha1.ComponentControlStatus{
							ValidAssessments: 1,
							TotalAssessments: 2,
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
			output := UpdateControls(testCase.inputComponentControlList, testCase.inputComponent)
			assert.Equal(t, output.Status.Controls, testCase.expectedOutput.Status.Controls)
			assert.Equal(t, output.Status.TotalControls, testCase.expectedOutput.Status.TotalControls)
			assert.Equal(t, output.Status.ImplementedControls, testCase.expectedOutput.Status.ImplementedControls)
		})
	}
}

func TestUpdateChild(t *testing.T) {
	metrics.SetUpMetrics()
	err := argusiov1alpha1.AddToScheme(scheme.Scheme)
	require.Nil(t, err)
	testCases := []struct {
		name           string
		expectedOutput string
		inputComponent *argusiov1alpha1.Component
		clientContent  []client.Object
	}{
		{
			name: "No Parents",
			inputComponent: &argusiov1alpha1.Component{
				Spec: argusiov1alpha1.ComponentSpec{
					Parents: []string{},
				},
			},
		},
		{
			name: "Parent Available",
			inputComponent: &argusiov1alpha1.Component{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: argusiov1alpha1.ComponentSpec{
					Parents: []string{"existing"},
				},
			},
			clientContent: []client.Object{
				&argusiov1alpha1.Component{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "existing",
						Namespace: "default",
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "Component",
						APIVersion: "argus.io/v1alpha1",
					},
				},
			},
		},
		{
			name: "Parent Not found",
			inputComponent: &argusiov1alpha1.Component{
				Spec: argusiov1alpha1.ComponentSpec{
					Parents: []string{"unexisting"},
				},
			},
			expectedOutput: "parent Component unexisting not found",
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			cl := fake.NewClientBuilder().WithObjects(testCase.clientContent...).WithStatusSubresource(testCase.clientContent...).Build()
			output := UpdateChild(context.TODO(), cl, testCase.inputComponent)
			if testCase.expectedOutput == "" {
				assert.Nil(t, output)
			} else {
				assert.ErrorContains(t, output, testCase.expectedOutput)
			}
		})
	}
}
