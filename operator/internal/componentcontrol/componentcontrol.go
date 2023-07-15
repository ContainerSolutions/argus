package componentcontrol

import (
	"context"
	"fmt"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/ContainerSolutions/argus/operator/internal/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetValidComponentAssessments(ctx context.Context, cl client.Client, res argusiov1alpha1.ComponentControl) ([]argusiov1alpha1.NamespacedName, int, error) {
	total := []argusiov1alpha1.NamespacedName{}
	valid := 0
	list := argusiov1alpha1.ComponentAssessmentList{}
	ComponentName, ok := res.Labels["argus.io/Component"]
	if !ok {
		return nil, 0, fmt.Errorf("object does not have expected label 'argus.io/Component'")
	}
	err := cl.List(ctx, &list, client.MatchingLabels{"argus.io/Component": ComponentName})
	if err != nil {
		return nil, 0, fmt.Errorf("could not list ComponentAssessment: %w", err)
	}
	for _, Assessment := range list.Items {
		if utils.Contains(res.Spec.RequiredAssessmentClasses, Assessment.Spec.Class) {
			if Assessment.Spec.ControlRef.Code == res.Spec.Definition.Code && Assessment.Spec.ControlRef.Version == res.Spec.Definition.Version {
				// This is a valid Assessment and should be in the list
				name := argusiov1alpha1.NamespacedName{
					Name:      Assessment.Name,
					Namespace: Assessment.Namespace,
				}
				total = append(total, name)
				if Assessment.Status.TotalAttestations == Assessment.Status.PassedAttestations && Assessment.Status.TotalAttestations > 0 {
					valid = valid + 1
				}
			}
		}
	}
	return total, valid, nil
}
