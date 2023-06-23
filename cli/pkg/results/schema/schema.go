package schema

import "argus/pkg/models"

type ResultDriver interface {
	Summary(c *models.Configuration)
	Detailed(c *models.Configuration)
	All(c *models.Configuration)
}

var Registry = make(map[string]ResultDriver)

func Register(name string, driver ResultDriver) {
	Registry[name] = driver
}
