package storage

import (
	_ "argus/pkg/storage/file"
	"argus/pkg/storage/schema"
	"fmt"
)

func Init(name string) (schema.StorageDriver, error) {
	driver, ok := schema.Registry[name]
	if !ok {
		return nil, fmt.Errorf("driver not found")
	}
	return driver, nil
}
