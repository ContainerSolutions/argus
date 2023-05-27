package attester

import (
	_ "argus/pkg/attester/command"
	"argus/pkg/attester/schema"
	"fmt"
)

func Init(name string) (schema.AttestDriver, error) {
	driver, ok := schema.Registry[name]
	if !ok {
		return nil, fmt.Errorf("driver not found")
	}
	return driver, nil
}
