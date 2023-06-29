package resourcerequirement

import (
	"context"
	"fmt"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/ContainerSolutions/argus/operator/internal/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetValidResourceImplementations(ctx context.Context, cl client.Client, res argusiov1alpha1.ResourceRequirement) ([]argusiov1alpha1.NamespacedName, int, error) {
	total := []argusiov1alpha1.NamespacedName{}
	valid := 0
	list := argusiov1alpha1.ResourceImplementationList{}
	resourceName, ok := res.Labels["argus.io/resource"]
	if !ok {
		return nil, 0, fmt.Errorf("object does not have expected label 'argus.io/resource'")
	}
	err := cl.List(ctx, &list, client.MatchingLabels{"argus.io/resource": resourceName})
	if err != nil {
		return nil, 0, fmt.Errorf("could not list ResourceImplementation: %w", err)
	}
	for _, implementation := range list.Items {
		if utils.Contains(res.Spec.RequiredImplementationClasses, implementation.Spec.Class) {
			if implementation.Spec.RequirementRef.Code == res.Spec.Definition.Code && implementation.Spec.RequirementRef.Version == res.Spec.Definition.Version {
				// This is a valid implementation and should be in the list
				name := argusiov1alpha1.NamespacedName{
					Name:      implementation.Name,
					Namespace: implementation.Namespace,
				}
				total = append(total, name)
				if implementation.Status.TotalAttestations == implementation.Status.PassedAttestations && implementation.Status.TotalAttestations > 0 {
					valid = valid + 1
				}
			}	
		}
	}
	return total, valid, nil
}
