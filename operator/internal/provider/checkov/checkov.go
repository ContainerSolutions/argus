package checkov

import (
	"fmt"
	"os"

	"os/exec"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	provider "github.com/ContainerSolutions/argus/operator/internal/provider/schema"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Client struct {
	RepoUrl string
	Checks  string
	Result  argusiov1alpha1.AttestationResultType
}

func (c *Client) Attest() (argusiov1alpha1.AttestationResult, error) {
	// ToDo: generate unique file names for clone_location and output_file_path
	clone_location := "/tmp/location"
	cmd := exec.Command("git", "clone", c.RepoUrl, clone_location)
	defer os.RemoveAll("/tmp/location")
	_, err := cmd.CombinedOutput()

	if err != nil {
		res := argusiov1alpha1.AttestationResult{
			Result: argusiov1alpha1.AttestationResultTypeUnknown,
			Logs:   err.Error(), // Logs has to be type string
			RunAt:  v1.Now(),
			Reason: fmt.Sprintf("could not get source repo for '%v'", c.RepoUrl),
		}
		return res, nil
	}

	checkov_cmd := exec.Command("checkov", "-d", clone_location, "--check", c.Checks, "-o", "cli")

	out, err := checkov_cmd.CombinedOutput()

	// Distinguish between execution and validation failure
	res := argusiov1alpha1.AttestationResult{
		Result: argusiov1alpha1.AttestationResultTypePass,
		Logs:   string(out),
		RunAt:  v1.Now(),
		Reason: "check pass",
	}
	if err != nil {
		res.Result = argusiov1alpha1.AttestationResultTypeUnknown
		res.Err = err.Error()
		res.Reason = "checkov execution returned error"
	}
	if checkov_cmd.ProcessState.ExitCode() != 0 {
		res.Result = argusiov1alpha1.AttestationResultTypeFail
		res.Reason = "checkov execution failed"
	}
	return res, nil
}

func (c *Client) Close() error {
	return nil
}

type Provider struct{}

func (p *Provider) New(spec *argusiov1alpha1.AttestationProviderSpec) (provider.AttestationClient, error) {
	c := &Client{}
	repourl_value, ok := spec.ProviderConfig["repo"]

	if ok {
		c.RepoUrl = repourl_value
	}

	checks_value, ok := spec.ProviderConfig["checks"]

	if ok {
		c.Checks = checks_value
	}

	return c, nil
}

func init() {
	provider.Register(&Provider{}, "checkov")
}
