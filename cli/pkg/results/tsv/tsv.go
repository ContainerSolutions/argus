package tsv

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ContainerSolutions/argus/cli/pkg/models"
	"github.com/ContainerSolutions/argus/cli/pkg/results/schema"
)

type TSVSummary struct {
}

func init() {
	schema.Register("tsv", &TSVSummary{})
}
func (t *TSVSummary) Summary(c *models.Configuration) {
	w := tabwriter.NewWriter(os.Stdout, 10, 4, 2, ' ', 0)
	var line string
	line = fmt.Sprintf("Resource\tStatus\tTotal Requirements\tImplemented Requirements\n")
	w.Write([]byte(line))
	for _, r := range c.Resources {
		line = fmt.Sprintf("%v\t%v\t%v\t%v\t\n", r.Name, r.Implemented, len(r.Requirements), r.ImplementedRequirements)
		w.Write([]byte(line))
	}
	w.Flush()
}

func (t *TSVSummary) All(c *models.Configuration) {

}

func (t *TSVSummary) Detailed(c *models.Configuration) {
	w := tabwriter.NewWriter(os.Stdout, 10, 4, 2, ' ', 0)
	var line string
	line = fmt.Sprintf("Resource\tRequirement\tImplementation\tAttestation\tEvaluated At\tResult\tLogs\n")
	w.Write([]byte(line))
	for _, r := range c.Resources {
		printReq := "N/A"
		printImp := "N/A"
		printAtt := "N/A"
		ranAt := "N/A"
		result := "N/A"
		logs := "N/A"
		if len(r.Requirements) == 0 {
			line = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t\n", r.Name, printReq, printImp, printAtt, ranAt, result, logs)
			w.Write([]byte(line))
		}
		for _, req := range r.Requirements {
			printReq = req.Requirement.Name
			if len(req.Implementations) == 0 {
				line = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t\n", r.Name, printReq, printImp, printAtt, ranAt, result, logs)
				w.Write([]byte(line))
			}
			for _, imp := range req.Implementations {
				printImp = imp.Implementation.Name
				if len(imp.Attestation) == 0 {
					line = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t\n", r.Name, printReq, printImp, printAtt, ranAt, result, logs)
					w.Write([]byte(line))
				}
				for _, ats := range imp.Attestation {
					printAtt = ats.Attestation.Name
					ranAt = fmt.Sprint(ats.Attestation.Result.RunAt)
					result = ats.Attestation.Result.Result
					logs = ats.Attestation.Result.Logs
					line = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t\n", r.Name, printReq, printImp, printAtt, ranAt, result, logs)
					w.Write([]byte(line))
				}
			}
		}
	}
	w.Flush()
}
