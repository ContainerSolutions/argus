package schema

import "argus/pkg/models"

type StorageDriver interface {
	Save(config *models.Configuration) error
	Load() (*models.Configuration, error)
	Configure(config map[string]interface{}) error
}

var Registry = make(map[string]StorageDriver)

func Register(name string, driver StorageDriver) {
	Registry[name] = driver
}
