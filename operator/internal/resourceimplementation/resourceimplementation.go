package resourceimplementation

import (
	"context"
	"fmt"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ListResourceAttestations(ctx context.Context, cl client.Client, res *argusiov1alpha1.ResourceImplementation) ([]argusiov1alpha1.ResourceAttestation, error) {
	list := argusiov1alpha1.ResourceAttestationList{}
	resourceName, ok := res.Labels["argus.io/resource"]
	if !ok {
		return nil, fmt.Errorf("object does not have expected label 'argus.io/resource'")
	}
	implementationName, ok := res.Labels["argus.io/implementation"]
	if !ok {
		return nil, fmt.Errorf("object does not have expected label 'argus.io/implementation'")
	}
	err := cl.List(ctx, &list, client.MatchingLabels{"argus.io/resource": resourceName, "argus.io/implementation": implementationName})
	if err != nil {
		return nil, fmt.Errorf("could not list ResourceImplementation: %w", err)
	}
	return list.Items, nil
}

func GetValidResourceAttestations(ctx context.Context, attestations []argusiov1alpha1.ResourceAttestation) ([]argusiov1alpha1.NamespacedName, int) {
	valid := 0
	children := []argusiov1alpha1.NamespacedName{}
	for _, attestation := range attestations {
		child := argusiov1alpha1.NamespacedName{
			Name:      attestation.Name,
			Namespace: attestation.Namespace,
		}
		children = append(children, child)
		if attestation.Status.Result.Result == argusiov1alpha1.AttestationResultTypePass {
			valid = valid + 1
		}
	}
	return children, valid
}
