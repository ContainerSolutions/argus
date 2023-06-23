package storage

import (
	"fmt"

	_ "github.com/ContainerSolutions/argus/cli/pkg/storage/file"
	"github.com/ContainerSolutions/argus/cli/pkg/storage/schema"
)

func Init(name string) (schema.StorageDriver, error) {
	driver, ok := schema.Registry[name]
	if !ok {
		return nil, fmt.Errorf("driver not found")
	}
	return driver, nil
}
