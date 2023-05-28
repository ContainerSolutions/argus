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

func (t *TSVSummary) Detailed(c *models.Configuration) {
	records := [][]string{
		{"Resource", "Requirement", "Implementation", "Attestation", "EvaluatedAt", "Result", "Logs"},
	}
	for _, r := range c.Resources {
		printReq := "N/A"
		printImp := "N/A"
		printAtt := "N/A"
		ranAt := "N/A"
		result := "N/A"
		logs := "N/A"
		if len(r.Requirements) == 0 {
			records = append(records, []string{r.Name, printReq, printImp, printAtt, ranAt, result, logs})
		}
		for _, req := range r.Requirements {
			printReq = req.Requirement.Name
			if len(req.Implementations) == 0 {
				records = append(records, []string{r.Name, printReq, printImp, printAtt, ranAt, result, logs})
			}
			for _, imp := range req.Implementations {
				printImp = imp.Implementaiton.Name
				if len(imp.Attestation) == 0 {
					records = append(records, []string{r.Name, printReq, printImp, printAtt, ranAt, result, logs})
				}
				for _, ats := range imp.Attestation {
					printAtt = ats.Attestation.Name
					ranAt = fmt.Sprint(ats.Attestation.Result.RunAt)
					result = ats.Attestation.Result.Result
					logs = ats.Attestation.Result.Logs
					records = append(records, []string{r.Name, printReq, printImp, printAtt, ranAt, result, logs})
				}
			}
		}
	}
	w := csv.NewWriter(os.Stdout)
	w.Comma = ';'
	w.WriteAll(records) // calls Flush internally
}
