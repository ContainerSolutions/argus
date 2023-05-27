package results

import (
	"argus/pkg/models"
	_ "argus/pkg/results/json"
	"argus/pkg/results/schema"
	_ "argus/pkg/results/tsv"
	"fmt"
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
