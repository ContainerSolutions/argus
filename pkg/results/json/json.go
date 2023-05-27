package json

import (
	"argus/pkg/models"
	"argus/pkg/results/schema"
	"encoding/json"
	"fmt"
	"os"
)

type JSONSummary struct {
}

func init() {
	schema.Register("json", &JSONSummary{})
}
func (t *JSONSummary) Summary(c *models.Configuration) {
	d := []struct {
		Resource                string
		Status                  string
		TotalRequirements       int
		ImplementedRequirements int
	}{}
	for _, r := range c.Resources {
		d = append(d, struct {
			Resource                string
			Status                  string
			TotalRequirements       int
			ImplementedRequirements int
		}{r.Name, fmt.Sprint(r.Implemented), len(r.Requirements), r.ImplementedRequirements})
	}
	w := json.NewEncoder(os.Stdout)
	w.Encode(d)
}
