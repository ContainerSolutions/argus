package resource

import (
	"context"
	"fmt"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func UpdateRequirements(resourceRequirementList argusiov1alpha1.ResourceRequirementList, resource *argusiov1alpha1.Resource) *argusiov1alpha1.Resource {
	validRequirements := 0
	reqs := make(map[string]*argusiov1alpha1.ResourceRequirementCompliance)
	for _, resourceRequirement := range resourceRequirementList.Items {
		status := argusiov1alpha1.ResourceRequirementCompliance{}
		status.Implemented = false
		if (resourceRequirement.Status.ValidImplementations == resourceRequirement.Status.TotalImplementations) || resourceRequirement.Status.TotalImplementations > 0 {
			status.Implemented = true
			validRequirements = validRequirements + 1
		}
		name := fmt.Sprintf("%v:%v", resourceRequirement.Spec.Definition.Code, resourceRequirement.Spec.Definition.Version)
		reqs[name] = &status
	}
	resource.Status.Requirements = reqs
	resource.Status.TotalRequirements = len(resourceRequirementList.Items)
	resource.Status.ImplementedRequirements = validRequirements
	return resource
}

func UpdateChild(ctx context.Context, cl client.Client, resource *argusiov1alpha1.Resource) error {
	var errorsA []error
	var allErrors error
	for _, parentName := range resource.Spec.Parents {
		parentResource := argusiov1alpha1.Resource{}
		namespacedName := types.NamespacedName{
			Name:      parentName,
			Namespace: resource.Namespace,
		}
		// TODO
		// Update Compliant status based on ResourceRequirement Status
		// TODO
		// Update parent adding current as a Child
		err := cl.Get(ctx, namespacedName, &parentResource)
		if err != nil {
			errorsA = append(errorsA, err)
		}
		original := parentResource.DeepCopy()
		if parentResource.Status.Children == nil {
			parentResource.Status.Children = make(map[string]argusiov1alpha1.ResourceChild)
		}
		parentResource.Status.Children[resource.Name] = argusiov1alpha1.ResourceChild{
			Compliant: false,
		}
		parentResource.Status.CompliantChildren = 1
		parentResource.Status.ImplementedRequirements = 1
		parentResource.Status.TotalRequirements = 1
		parentResource.Status.TotalChildren = 1
		err = cl.Status().Patch(ctx, &parentResource, client.MergeFrom(original))
		if err != nil {
			errorsA = append(errorsA, err)
		}
	}
	if errorsA != nil {
		for _, err := range errorsA {
			allErrors = errors.Wrap(allErrors, fmt.Sprintf("error occurred: %v", err.Error()))
		}

		return fmt.Errorf("could not update parent resources: %v", errorsA)
	}
	return allErrors
}
