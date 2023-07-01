package fake

import (
	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	provider "github.com/ContainerSolutions/argus/operator/internal/provider/schema"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Client struct {
}

func (c *Client) Attest() (argusiov1alpha1.AttestationResult, error) {
	res := argusiov1alpha1.AttestationResult{
		Result: argusiov1alpha1.AttestationResultTypePass,
		Logs:   "fake",
		RunAt:  v1.Now(),
		Reason: "fake provider always pass",
	}
	return res, nil
}

func (c *Client) Close() error {
	return nil
}

type Provider struct{}

func (p *Provider) New(spec *argusiov1alpha1.AttestationProviderSpec) (provider.AttestationClient, error) {
	c := &Client{}
	return c, nil
}

func init() {
	provider.Register(&Provider{}, "fake")
}
