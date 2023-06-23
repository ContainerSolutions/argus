package command

import (
	"errors"
	"github.com/ContainerSolutions/argus/cli/pkg/models"
	"testing"

	"gotest.tools/v3/assert"
)

type FakeCommand struct {
	CombinedOutputFn func() ([]byte, error)
	StringFn         func() string
	ExitCodeFn       func() int
}

func (f *FakeCommand) Command(name string, arg ...string) {
}

func (f *FakeCommand) CombinedOutput() ([]byte, error) {
	return f.CombinedOutputFn()
}

func (f *FakeCommand) String() string {
	return f.StringFn()
}
func (f *FakeCommand) ExitCode() int {
	return f.ExitCodeFn()
}

func WithCombinedOutput(data []byte, err error) func() ([]byte, error) {
	return func() ([]byte, error) {
		return data, err
	}
}

func WithString(data string) func() string {
	return func() string {
		return data
	}
}

func WithExitCode(data int) func() int {
	return func() int {
		return data
	}
}

func TestCommand(t *testing.T) {
	f := &FakeCommand{}
	a := AttestCommand{
		cmd: f,
	}
	testCase := []struct {
		name             string
		combinedOutputFn func() ([]byte, error)
		stringFn         func() string
		exitCodefn       func() int
		expectedErr      string
		input            *models.Attestation
		expectedResult   *models.AttestationResult
	}{
		{
			name:             "NoOutputExitCodeMatches",
			combinedOutputFn: WithCombinedOutput([]byte("test"), nil),
			stringFn:         WithString("fake"),
			exitCodefn:       WithExitCode(0),
			input: &models.Attestation{
				Name: "fake",
				Type: "fake",
				CommandRef: models.AttestationByCommand{
					Command:          "fake",
					Args:             []string{},
					ExpectedExitCode: 0,
					ExpectedOutput:   "",
				},
			},
			expectedResult: &models.AttestationResult{
				Command: "fake",
				Logs:    "$ fake:\ntest",
				Result:  "PASS",
				Reason:  "",
				Err:     "",
			},
		},
		{
			name:             "CodeFail",
			combinedOutputFn: WithCombinedOutput([]byte("test"), nil),
			stringFn:         WithString("fake"),
			exitCodefn:       WithExitCode(1),
			input: &models.Attestation{
				Name: "fake",
				Type: "fake",
				CommandRef: models.AttestationByCommand{
					Command:          "fake",
					Args:             []string{},
					ExpectedExitCode: 0,
					ExpectedOutput:   "",
				},
			},
			expectedResult: &models.AttestationResult{
				Command: "fake",
				Logs:    "$ fake:\ntest",
				Result:  "FAIL",
				Reason:  "Code failed! Got 1 But Expected 0\n",
				Err:     "",
			},
		},
		{
			name:             "OutputFail",
			combinedOutputFn: WithCombinedOutput([]byte("fail"), nil),
			stringFn:         WithString("fake"),
			exitCodefn:       WithExitCode(0),
			input: &models.Attestation{
				Name: "fake",
				Type: "fake",
				CommandRef: models.AttestationByCommand{
					Command:          "fake",
					Args:             []string{},
					ExpectedExitCode: 0,
					ExpectedOutput:   "pass",
				},
			},
			expectedResult: &models.AttestationResult{
				Command: "fake",
				Logs:    "$ fake:\nfail",
				Result:  "FAIL",
				Reason:  "Output Check failed! Wanted 'pass'\nGot 'fail'\n",
				Err:     "",
			},
		},
		{
			name:             "CommandError",
			combinedOutputFn: WithCombinedOutput([]byte("fail"), errors.New("with error")),
			stringFn:         WithString("fake"),
			exitCodefn:       WithExitCode(0),
			input: &models.Attestation{
				Name: "fake",
				Type: "fake",
				CommandRef: models.AttestationByCommand{
					Command:          "fake",
					Args:             []string{},
					ExpectedExitCode: 0,
					ExpectedOutput:   "pass",
				},
			},
			expectedResult: &models.AttestationResult{
				Command: "fake",
				Logs:    "$ fake:\nfail",
				Result:  "FAIL",
				Reason:  "Output Check failed! Wanted 'pass'\nGot 'fail'\n",
				Err:     "with error",
			},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			f.CombinedOutputFn = tc.combinedOutputFn
			f.StringFn = tc.stringFn
			f.ExitCodeFn = tc.exitCodefn
			res, err := a.Attest(tc.input)
			if tc.expectedErr != "" {
				assert.ErrorContains(t, err, tc.expectedErr)
			}
			if tc.expectedResult != nil {
				assert.Equal(t, tc.expectedResult.Command, res.Command)
				assert.Equal(t, tc.expectedResult.Err, res.Err)
				assert.Equal(t, tc.expectedResult.Logs, res.Logs)
				assert.Equal(t, tc.expectedResult.Reason, res.Reason)
				assert.Equal(t, tc.expectedResult.Result, res.Result)
			}
		})
	}
}
