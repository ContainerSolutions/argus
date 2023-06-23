package attester

import (
	"fmt"

	_ "github.com/ContainerSolutions/argus/cli/pkg/attester/command"
	"github.com/ContainerSolutions/argus/cli/pkg/attester/schema"
)

func Init(name string) (schema.AttestDriver, error) {
	driver, ok := schema.Registry[name]
	if !ok {
		return nil, fmt.Errorf("driver not found")
	}
	return driver, nil
}
