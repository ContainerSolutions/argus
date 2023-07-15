package fake

import (
	"fmt"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	provider "github.com/ContainerSolutions/argus/operator/internal/provider/schema"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Client struct {
	Result argusiov1alpha1.AttestationResultType
}

func (c *Client) Attest() (argusiov1alpha1.AttestationResult, error) {
	res := argusiov1alpha1.AttestationResult{
		Result: c.Result,
		Logs:   "fake",
		RunAt:  v1.Now(),
		Reason: fmt.Sprintf("fake provider configured for '%v'", c.Result),
	}
	return res, nil
}

func (c *Client) Close() error {
	return nil
}

type Provider struct{}

func (p *Provider) New(name string, spec *argusiov1alpha1.AttestationProviderSpec) (provider.AttestationClient, error) {
	var res argusiov1alpha1.AttestationResultType
	r, ok := spec.ProviderConfig["result"]
	if !ok {
		return &Client{Result: argusiov1alpha1.AttestationResultTypePass}, nil
	}
	switch r {
	case "Pass":
		res = argusiov1alpha1.AttestationResultTypePass
	case "Fail":
		res = argusiov1alpha1.AttestationResultTypeFail
	default:
		res = argusiov1alpha1.AttestationResultTypeUnknown
	}
	c := &Client{
		Result: res,
	}
	return c, nil
}

func init() {
	provider.Register(&Provider{}, "fake")
}
