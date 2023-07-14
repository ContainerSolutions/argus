package checkov

import (
	"fmt"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	provider "github.com/ContainerSolutions/argus/operator/internal/provider/schema"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os/exec"
)

type Client struct {
	RepoUrl string
	Checks  string
	Result  argusiov1alpha1.AttestationResultType
}

func (c *Client) Attest() (argusiov1alpha1.AttestationResult, error) {
	// ToDo: generate unique file names for clone_location and output_file_path
	clone_location := "/tmp/location"
	output_file_path := "/tmp/output"

	var result_type argusiov1alpha1.AttestationResultType

	cmd := exec.Command("git", "clone", c.RepoUrl, clone_location)
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
		result_type = argusiov1alpha1.AttestationResultTypeFail

		res := argusiov1alpha1.AttestationResult{
			Result: result_type,
			Logs:   err.Error(), // Logs has to be type string
			RunAt:  v1.Now(),
			Reason: fmt.Sprintf("checkov provider configured for '%v'", result_type),
		}
		// ToDo: delete clone_location
		return res, nil
	}

	// checkov command
	checkov_cmd := exec.Command("checkov", "-d", clone_location, "--check", c.Checks, "--output-file-path", output_file_path)

	err = checkov_cmd.Run()

	// Distinguish between execution and validation failure
	if err != nil {
		fmt.Println(err)
		result_type = argusiov1alpha1.AttestationResultTypeFail

		res := argusiov1alpha1.AttestationResult{
			Result: result_type,
			Logs:   err.Error(), // Logs has to be type string
			RunAt:  v1.Now(),
			Reason: fmt.Sprintf("checkov provider configured for '%v'", result_type),
		}
		// ToDo: delete output_file
		return res, nil
	} else {
		result_type = argusiov1alpha1.AttestationResultTypePass
	}

	content, err := ioutil.ReadFile(output_file_path) // content is of type []byte

	if err != nil {
		result_type = argusiov1alpha1.AttestationResultTypeFail
		content = []byte(err.Error()) // Converted err to type []byte
	}

	res := argusiov1alpha1.AttestationResult{
		Result: result_type,
		Logs:   string(content), // Logs has to be type string
		RunAt:  v1.Now(),
		Reason: fmt.Sprintf("checkov provider configured for '%v'", result_type),
	}
	// ToDo: delete clone_location
	return res, nil
}

func (c *Client) Close() error {
	return nil
}

type Provider struct{}

func (p *Provider) New(spec *argusiov1alpha1.AttestationProviderSpec) (provider.AttestationClient, error) {
	c := &Client{}
	repourl_value, ok := spec.ProviderConfig["repourl"]

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
