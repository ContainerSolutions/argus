package componentassessment

import (
	"context"
	"fmt"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ListComponentAttestations(ctx context.Context, cl client.Client, res *argusiov1alpha1.ComponentAssessment) ([]argusiov1alpha1.ComponentAttestation, error) {
	list := argusiov1alpha1.ComponentAttestationList{}
	ComponentName, ok := res.Labels["argus.io/Component"]
	if !ok {
		return nil, fmt.Errorf("object does not have expected label 'argus.io/Component'")
	}
	AssessmentName, ok := res.Labels["argus.io/Assessment"]
	if !ok {
		return nil, fmt.Errorf("object does not have expected label 'argus.io/Assessment'")
	}
	err := cl.List(ctx, &list, client.MatchingLabels{"argus.io/Component": ComponentName, "argus.io/Assessment": AssessmentName})
	if err != nil {
		return nil, fmt.Errorf("could not list ComponentAssessment: %w", err)
	}
	return list.Items, nil
}

func GetValidComponentAttestations(ctx context.Context, attestations []argusiov1alpha1.ComponentAttestation) ([]argusiov1alpha1.NamespacedName, int) {
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
