package requirement

import (
	"context"
	"fmt"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/ContainerSolutions/argus/operator/internal/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetResourceRequirementsFromRequirement(ctx context.Context, cl client.Client, requirement *argusiov1alpha1.Requirement) (map[string]argusiov1alpha1.ResourceRequirement, error) {
	resourceRequirementList := argusiov1alpha1.ResourceRequirementList{}
	err := cl.List(ctx, &resourceRequirementList, client.MatchingLabels{"argus.io/requirement": fmt.Sprintf("%v_%v", requirement.Spec.Definition.Code, requirement.Spec.Definition.Version)})
	if err != nil {
		return nil, fmt.Errorf("could not list ResourceRequirements: %w", err)
	}
	resReqs := make(map[string]argusiov1alpha1.ResourceRequirement)
	for _, resourceRequirement := range resourceRequirementList.Items {
		resReqs[resourceRequirement.Name] = resourceRequirement
	}
	return resReqs, nil
}

func LifecycleResourceRequirements(ctx context.Context, cl client.Client, classes []string, resources []argusiov1alpha1.Resource, resReq map[string]argusiov1alpha1.ResourceRequirement) error {
	resourceNames := []string{}
	for _, resource := range resources {
		resourceNames = append(resourceNames, resource.Name)
	}
	for _, resourceRequirement := range resReq {
		refResource, ok := resourceRequirement.ObjectMeta.Labels["argus.io/resource"]
		if !ok {
			return fmt.Errorf("object '%v' does not contain expected label 'argus.io/resource'", resourceRequirement.Name)
		}
		// If resource does not exist, it was deleted - we need to delete resourceRequirement
		if !utils.Contains(resourceNames, refResource) {
			err := cl.Delete(ctx, &resourceRequirement)
			if err != nil {
				return fmt.Errorf("could not delete ResourceRequirement '%v': %w", resourceRequirement.Name, err)
			}
		}
		class, ok := resourceRequirement.ObjectMeta.Labels["argus.io/resource-class"]
		if !ok {
			return fmt.Errorf("object '%v' does not contain expected label 'argus.io/resource-class'", resourceRequirement.Name)
		}
		// If requirement Class has changed, we need to delete resourceRequirement
		if !utils.Contains(classes, class) {
			err := cl.Delete(ctx, &resourceRequirement)
			if err != nil {
				return fmt.Errorf("could not delete ResourceRequirement '%v': %w", resourceRequirement.Name, err)
			}
		}
		// If resource Class has changed, we need to delete ResourceRequirement
		for _, resource := range resources {
			refResource, ok := resourceRequirement.ObjectMeta.Labels["argus.io/resource"]
			if !ok {
				return fmt.Errorf("object '%v' does not contain expected label 'argus.io/resource'", resourceRequirement.Name)
			}
			if refResource == resource.Name && !utils.Contains(resource.Spec.Classes, class) {
				err := cl.Delete(ctx, &resourceRequirement)
				if err != nil {
					return fmt.Errorf("could not delete ResourceRequirement '%v': %w", resourceRequirement.Name, err)
				}
			}
		}
	}
	return nil
}

func CreateOrUpdateResourceRequirements(ctx context.Context, cl client.Client, req *argusiov1alpha1.Requirement, resources []argusiov1alpha1.Resource) error {
	for _, class := range req.Spec.ApplicableResourceClasses {
		for _, resource := range resources {
			if utils.Contains(resource.Spec.Classes, class) {
				resReq := &argusiov1alpha1.ResourceRequirement{
					ObjectMeta: metav1.ObjectMeta{
						Name:      fmt.Sprintf("%v-%v", req.Name, resource.Name),
						Namespace: req.Namespace,
					},
				}
				emptyMutation := func() error {
					resReq.Spec.Definition = req.Spec.Definition
					resReq.Spec.RequiredImplementationClasses = req.Spec.RequiredImplementationClasses
					resReq.ObjectMeta.Labels = map[string]string{
						"argus.io/requirement":    fmt.Sprintf("%v_%v", req.Spec.Definition.Code, req.Spec.Definition.Version),
						"argus.io/resource":       resource.Name,
						"argus.io/resource-class": class,
					}
					return nil
				}
				_, err := ctrl.CreateOrUpdate(ctx, cl, resReq, emptyMutation)
				if err != nil {
					return fmt.Errorf("could not create resourcerequirement '%v': %w", resReq.Name, err)
				}
			}
		}
	}
	return nil
}
