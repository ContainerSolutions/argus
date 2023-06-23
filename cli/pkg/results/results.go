package results

import (
	"fmt"

	"github.com/ContainerSolutions/argus/cli/pkg/models"
	_ "github.com/ContainerSolutions/argus/cli/pkg/results/json"
	"github.com/ContainerSolutions/argus/cli/pkg/results/schema"
	_ "github.com/ContainerSolutions/argus/cli/pkg/results/tsv"
)

func Summary(config *models.Configuration, format string) error {
	sum, ok := schema.Registry[format]
	if !ok {
		return fmt.Errorf("summary format %v is not supported", format)
	}
	sum.Summary(config)
	return nil
}

func Detailed(config *models.Configuration, format string) error {
	sum, ok := schema.Registry[format]
	if !ok {
		return fmt.Errorf("summary format %v is not supported", format)
	}
	sum.Detailed(config)
	return nil
}

func All(config *models.Configuration, format string) error {
	sum, ok := schema.Registry[format]
	if !ok {
		return fmt.Errorf("summary format %v is not supported", format)
	}
	sum.All(config)
	return nil
}
