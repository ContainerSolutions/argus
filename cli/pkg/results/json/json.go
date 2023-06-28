package json

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ContainerSolutions/argus/cli/pkg/models"
	"github.com/ContainerSolutions/argus/cli/pkg/results/schema"
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
	err := w.Encode(d)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error happened encoding:%v", err)
	}
}

func (t *JSONSummary) Detailed(c *models.Configuration) {
	d := []struct {
		Resource       string
		Implementation string
		Requirement    string
		Attestation    string
		EvaluatedAt    time.Time
		Result         string
		Logs           string
	}{}
	for _, r := range c.Resources {
		jimp := ""
		jreq := ""
		att := ""
		eval := time.Time{}
		res := ""
		logs := ""
		if len(r.Requirements) == 0 {
			d = append(d, struct {
				Resource       string
				Implementation string
				Requirement    string
				Attestation    string
				EvaluatedAt    time.Time
				Result         string
				Logs           string
			}{Resource: r.Name, Requirement: jreq, Implementation: jimp, Attestation: att, EvaluatedAt: eval, Result: res, Logs: logs})
		}
		for _, req := range r.Requirements {
			jreq = req.Requirement.Name
			if len(req.Implementations) == 0 {
				d = append(d, struct {
					Resource       string
					Implementation string
					Requirement    string
					Attestation    string
					EvaluatedAt    time.Time
					Result         string
					Logs           string
				}{Resource: r.Name, Requirement: jreq, Implementation: jimp, Attestation: att, EvaluatedAt: eval, Result: res, Logs: logs})
			}
			for _, imp := range req.Implementations {
				jimp = imp.Implementation.Name
				if len(imp.Attestation) == 0 {
					d = append(d, struct {
						Resource       string
						Implementation string
						Requirement    string
						Attestation    string
						EvaluatedAt    time.Time
						Result         string
						Logs           string
					}{Resource: r.Name, Requirement: jreq, Implementation: jimp, Attestation: att, EvaluatedAt: eval, Result: res, Logs: logs})
				}
				for _, ats := range imp.Attestation {
					att = ats.Attestation.Name
					eval = ats.Attestation.Result.RunAt
					res = ats.Attestation.Result.Result
					logs = ats.Attestation.Result.Logs
					d = append(d, struct {
						Resource       string
						Implementation string
						Requirement    string
						Attestation    string
						EvaluatedAt    time.Time
						Result         string
						Logs           string
					}{Resource: r.Name, Requirement: jreq, Implementation: jimp, Attestation: att, EvaluatedAt: eval, Result: res, Logs: logs})
				}
			}
		}

	}
	w := json.NewEncoder(os.Stdout)
	err := w.Encode(d)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error happened while encoding:%v", err)
	}
}

func (t *JSONSummary) All(c *models.Configuration) {
	w := json.NewEncoder(os.Stdout)
	err := w.Encode(c.Resources)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error happened while encoding:%v", err)
	}
}
