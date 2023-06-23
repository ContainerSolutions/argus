package schema

import "github.com/ContainerSolutions/argus/cli/pkg/models"

type AttestDriver interface {
	Attest(c *models.Attestation) (*models.AttestationResult, error)
}

var Registry = make(map[string]AttestDriver)

func Register(name string, driver AttestDriver) {
	Registry[name] = driver
}
