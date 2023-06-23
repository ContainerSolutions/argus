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
	cmd Command
}

type Command interface {
	Command(name string, arg ...string)
	CombinedOutput() ([]byte, error)
	String() string
	ExitCode() int
}
type ExecCommand struct {
	command *exec.Cmd
}

func (e *ExecCommand) Command(name string, arg ...string) {
	e.command = exec.Command(name, arg...)
}

func (e *ExecCommand) CombinedOutput() ([]byte, error) {
	return e.command.CombinedOutput()
}

func (e *ExecCommand) String() string {
	return e.command.String()
}
func (e *ExecCommand) ExitCode() int {
	return e.command.ProcessState.ExitCode()
}

func init() {
	schema.Register("command", &AttestCommand{
		cmd: &ExecCommand{},
	})
}

func (f *AttestCommand) Attest(a *models.Attestation) (*models.AttestationResult, error) {
	f.cmd.Command(a.CommandRef.Command, a.CommandRef.Args...)
	res := models.AttestationResult{}
	out, err := f.cmd.CombinedOutput()
	res.Command = f.cmd.String()
	res.Logs = fmt.Sprintf("$ %v:\n%v", res.Command, string(out))
	res.RunAt = time.Now()
	if err != nil {
		res.Err = err.Error()
	}
	if f.cmd.ExitCode() != a.CommandRef.ExpectedExitCode {
		res.Result = "FAIL"
		res.Reason = fmt.Sprintf("Code failed! Got %v But Expected %v\n", f.cmd.ExitCode(), a.CommandRef.ExpectedExitCode)
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
