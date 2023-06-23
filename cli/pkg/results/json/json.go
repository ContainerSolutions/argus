package json

import (
	"argus/pkg/models"
	"argus/pkg/results/schema"
	"encoding/json"
	"fmt"
	"os"
	"time"
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
	w.Encode(d)
}

func (t *JSONSummary) All(c *models.Configuration) {
	w := json.NewEncoder(os.Stdout)
	w.Encode(c.Resources)
}
