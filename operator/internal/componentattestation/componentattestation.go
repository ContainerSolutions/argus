package componentattestation

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/ContainerSolutions/argus/operator/internal/provider"
	"github.com/ContainerSolutions/argus/operator/internal/provider/schema"
)

func GetAttestationClient(ctx context.Context, cl client.Client, res *argusiov1alpha1.ComponentAttestation) (schema.AttestationClient, error) {
	providerSpec := argusiov1alpha1.AttestationProvider{}
	req := types.NamespacedName{
		Name:      res.Spec.ProviderRef.Name,
		Namespace: res.Spec.ProviderRef.Namespace,
	}
	err := cl.Get(ctx, req, &providerSpec)
	if err != nil {
		return nil, fmt.Errorf("could not get provider spec '%v': %w", req.Name, err)
	}
	prov, err := provider.GetProvider(providerSpec.Spec.Type)
	if err != nil {
		return nil, fmt.Errorf("could not get provider '%v': %w", req.Name, err)
	}
	attestationClient, err := prov.New(&providerSpec.Spec)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate client for provider '%v': %w", req.Name, err)
	}
	return attestationClient, nil
}
