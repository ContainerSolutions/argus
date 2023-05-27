package command

import (
	"argus/pkg/attester/schema"
	"argus/pkg/models"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type AttestCommand struct {
}

func init() {
	schema.Register("command", &AttestCommand{})
}

func (f *AttestCommand) Attest(a *models.Attestation) (*models.AttestationResult, error) {
	cmd := exec.Command(a.CommandRef.Command, a.CommandRef.Args...)
	res := models.AttestationResult{}
	out, err := cmd.CombinedOutput()
	res.Command = cmd.String()
	res.Logs = fmt.Sprintf("$ %v:\n%v", res.Command, string(out))
	res.RunAt = time.Now()
	if err != nil {
		res.Err = err.Error()
	}
	if cmd.ProcessState.ExitCode() != a.CommandRef.ExpectedExitCode {
		res.Result = "FAIL"
		res.Reason = fmt.Sprintf("Code failed! Got %v But Expected %v\n", cmd.ProcessState.ExitCode(), a.CommandRef.ExpectedExitCode)
		a.Result = res
		return &res, nil
	}
	if a.CommandRef.ExpectedOutput != "" {
		if !strings.Contains(string(out), a.CommandRef.ExpectedOutput) {
			res.Result = "FAIL"
			res.Reason = fmt.Sprintf("Output Check failed! Wanted '%v'\nGot '%v'\n", a.CommandRef.ExpectedOutput, string(out))
			a.Result = res
			return &res, nil
		}
	}
	res.Result = "PASS"
	res.Reason = ""
	a.Result = res
	return &res, nil
}
