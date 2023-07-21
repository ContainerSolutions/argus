package command

import (
	"os/exec"
	"strconv"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	provider "github.com/ContainerSolutions/argus/operator/internal/provider/schema"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Client struct {
	Command            string
	ExpectedStatusCode int
}

func (c *Client) Attest() (argusiov1alpha1.AttestationResult, error) {
	cmd := exec.Command(c.Command)
	out, err := cmd.CombinedOutput()
	result := argusiov1alpha1.AttestationResultTypePass
	if cmd.ProcessState.ExitCode() != c.ExpectedStatusCode {
		result = argusiov1alpha1.AttestationResultTypeFail
	}
	res := argusiov1alpha1.AttestationResult{
		Result: result,
		Logs:   string(out),
		RunAt:  v1.Now(),
		Reason: "command execution output",
	}
	if err != nil {
		res.Err = err.Error()
	}
	return res, nil
}

func (c *Client) Close() error {
	return nil
}

type Provider struct{}

func (p *Provider) New(name string, spec *argusiov1alpha1.AttestationProviderSpec) (provider.AttestationClient, error) {

	c := &Client{}
	cmdString, ok := spec.ProviderConfig["cmd"]
	if ok {
		c.Command = cmdString
	}
	statusCode, ok := spec.ProviderConfig["expectedStatusCode"]
	if ok {
		c.ExpectedStatusCode, _ = strconv.Atoi(statusCode) //nolint
	}
	return c, nil
}

func init() {
	provider.Register(&Provider{}, "command")
}
