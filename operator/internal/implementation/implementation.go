package implementation

import (
	"context"
	"fmt"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/ContainerSolutions/argus/operator/internal/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetResourceImplementations(ctx context.Context, cl client.Client, res *argusiov1alpha1.Implementation) (map[string]argusiov1alpha1.ResourceImplementation, error) {
	resourceImplementationList := argusiov1alpha1.ResourceImplementationList{}
	err := cl.List(ctx, &resourceImplementationList, client.MatchingLabels{"argus.io/implementation": res.Name})
	if err != nil {
		return nil, fmt.Errorf("could not list ResourceImplementation: %w", err)
	}
	resReqs := make(map[string]argusiov1alpha1.ResourceImplementation)
	for _, item := range resourceImplementationList.Items {
		resReqs[item.Name] = item
	}
	return resReqs, nil
}

func LifecycleResourceImplementations(ctx context.Context, cl client.Client, resources []argusiov1alpha1.Resource, items map[string]argusiov1alpha1.ResourceImplementation) error {
	resourceNames := []string{}
	for _, resource := range resources {
		resourceNames = append(resourceNames, resource.Name)
	}
	for _, item := range items {
		refResource, ok := item.ObjectMeta.Labels["argus.io/resource"]
		if !ok {
			return fmt.Errorf("object '%v' does not contain expected label 'argus.io/resource'", item.Name)
		}
		// If resource does not exist, it was deleted - we need to delete resourceRequirement
		if !utils.Contains(resourceNames, refResource) {
			err := cl.Delete(ctx, &item)
			if err != nil {
				return fmt.Errorf("could not delete ResourceImplementation '%v': %w", item.Name, err)
			}
		}
	}
	return nil
}

func CreateOrUpdateResourceImplementations(ctx context.Context, cl client.Client, res *argusiov1alpha1.Implementation, resources []argusiov1alpha1.Resource) ([]argusiov1alpha1.NamespacedName, error) {
	all := []argusiov1alpha1.NamespacedName{}
	// Treat Cascading policy. In order to do that, we need to add every child which this implementation targets.
	resourceNameList := []string{}
	resourceMap := make(map[string]argusiov1alpha1.Resource)
	for _, resource := range resources {
		resourceMap[resource.Name] = resource
	}
	for _, resource := range resources {
		for _, refs := range res.Spec.ResourceRef {
			// If this resource is actually targetted by this implementation
			if refs.Name == resource.Name && refs.Namespace == resource.Namespace {
				resourceNameList = append(resourceNameList, resource.Name)
				if res.Spec.CascadePolicy == argusiov1alpha1.CascadingPolicyCascade {
					for childName := range resource.Status.Children {
						// Add children
						resourceNameList = append(resourceNameList, childName)
					}
				}
			}

		}
	}
	for _, resourceName := range resourceNameList {
		resource := resourceMap[resourceName]
		resImp := &argusiov1alpha1.ResourceImplementation{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%v-%v", res.Name, resource.Name),
				Namespace: res.Namespace,
			},
		}
		emptyMutation := func() error {
			resImp.Spec.RequirementRef = res.Spec.RequirementRef
			resImp.Spec.Class = res.Spec.Class
			resImp.ObjectMeta.Labels = map[string]string{
				"argus.io/implementation": res.Name,
				"argus.io/resource":       resource.Name,
				"argus.io/requirement":    fmt.Sprintf("%v_%v", res.Spec.RequirementRef.Code, res.Spec.RequirementRef.Version),
			}
			return nil
		}
		_, err := ctrl.CreateOrUpdate(ctx, cl, resImp, emptyMutation)
		if err != nil {
			return nil, fmt.Errorf("could not create resourceImplementation '%v': %w", resImp.Name, err)
		}
		child := argusiov1alpha1.NamespacedName{
			Name:      resImp.Name,
			Namespace: resImp.Namespace,
		}
		all = append(all, child)
	}
	return all, nil
}
