package tsv

import (
	"argus/pkg/models"
	"argus/pkg/results/schema"
	"encoding/csv"
	"fmt"
	"os"
)

type TSVSummary struct {
}

func init() {
	schema.Register("tsv", &TSVSummary{})
}
func (t *TSVSummary) Summary(c *models.Configuration) {
	records := [][]string{
		{"Resource", "Status", "Total Requirements", "Implemented Requirements"},
	}
	for _, r := range c.Resources {
		records = append(records, []string{r.Name, fmt.Sprint(r.Implemented), fmt.Sprint(len(r.Requirements)), fmt.Sprint(r.ImplementedRequirements)})
	}
	w := csv.NewWriter(os.Stdout)
	w.Comma = '\t'
	w.WriteAll(records) // calls Flush internally
}
