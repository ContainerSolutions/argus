package attestation

import (
	"context"
	"fmt"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/ContainerSolutions/argus/operator/internal/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func GetResourceAttestations(ctx context.Context, cl client.Client, res *argusiov1alpha1.Attestation) (map[string]argusiov1alpha1.ResourceAttestation, error) {
	resourceAttestationList := argusiov1alpha1.ResourceAttestationList{}
	err := cl.List(ctx, &resourceAttestationList, client.MatchingLabels{"argus.io/attestation": res.Name})
	if err != nil {
		return nil, fmt.Errorf("could not list ResourceAttestation: %w", err)
	}
	resReqs := make(map[string]argusiov1alpha1.ResourceAttestation)
	for _, item := range resourceAttestationList.Items {
		resReqs[item.Name] = item
	}
	return resReqs, nil
}

func LifecycleResourceAttestations(ctx context.Context, cl client.Client, implementationRef string, resources []argusiov1alpha1.ResourceImplementation, items map[string]argusiov1alpha1.ResourceAttestation) error {
	resourceNames := []string{}
	for _, resource := range resources {
		label, ok := resource.ObjectMeta.Labels["argus.io/resource"]
		if !ok {
			return fmt.Errorf("resource implementation '%v' does not contain expected label 'argus.io/resource'", resource.Name)
		}
		resourceNames = append(resourceNames, label)
	}
	for _, item := range items {
		// if item does not belong anymore to the same implementation ref, delete it (as the attestation was updated)
		if item.Labels["argus.io/implementation"] != implementationRef {
			err := cl.Delete(ctx, &item)
			if err != nil {
				return fmt.Errorf("could not delete ResourceAttestation '%v': %w", item.Name, err)
			}
			continue
		}
		refResource, ok := item.ObjectMeta.Labels["argus.io/resource"]
		if !ok {
			return fmt.Errorf("object '%v' does not contain expected label 'argus.io/resource'", item.Name)
		}
		// If resource does not exist, it was deleted - we need to delete resourceRequirement
		if !utils.Contains(resourceNames, refResource) {
			err := cl.Delete(ctx, &item)
			if err != nil {
				return fmt.Errorf("could not delete ResourceAttestation '%v': %w", item.Name, err)
			}
		}
	}
	return nil

}

func CreateOrUpdateResourceAttestations(ctx context.Context, cl client.Client, scheme *runtime.Scheme, res *argusiov1alpha1.Attestation, resources []argusiov1alpha1.ResourceImplementation) ([]argusiov1alpha1.NamespacedName, error) {
	all := []argusiov1alpha1.NamespacedName{}
	for _, resource := range resources {
		implementationName, ok := resource.ObjectMeta.Labels["argus.io/implementation"]
		if !ok {
			return nil, fmt.Errorf("resource implementation '%v' does not contain expected label 'argus.io/implementation'", resource.Name)
		}
		if res.Spec.ImplementationRef == implementationName {
			resourceName, ok := resource.ObjectMeta.Labels["argus.io/resource"]
			if !ok {
				return nil, fmt.Errorf("resource implementation '%v' does not contain expected label 'argus.io/resource'", resource.Name)
			}
			requirementName, ok := resource.ObjectMeta.Labels["argus.io/requirement"]
			if !ok {
				return nil, fmt.Errorf("resource implementation '%v' does not contain expected label 'argus.io/requirement'", resource.Name)
			}
			resAtt := &argusiov1alpha1.ResourceAttestation{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%v-%v", res.Name, resourceName),
					Namespace: res.Namespace,
				},
			}
			controllerutil.SetControllerReference(res, &resAtt.ObjectMeta, scheme)
			emptyMutation := func() error {
				resAtt.Spec.ProviderRef = res.Spec.ProviderRef
				resAtt.ObjectMeta.Labels = map[string]string{
					"argus.io/implementation": res.Spec.ImplementationRef,
					"argus.io/attestation":    res.Name,
					"argus.io/resource":       resourceName,
					"argus.io/requirement":    requirementName,
				}
				return nil
			}
			_, err := ctrl.CreateOrUpdate(ctx, cl, resAtt, emptyMutation)
			if err != nil {
				return nil, fmt.Errorf("could not create ResourceAttestation '%v': %w", resAtt.Name, err)
			}
			child := argusiov1alpha1.NamespacedName{
				Name:      resAtt.Name,
				Namespace: resAtt.Namespace,
			}
			all = append(all, child)
		}
	}
	return all, nil
}
