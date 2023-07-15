package file

// This provider parses any http page and returns a regular expression match (and how many times).
// It allows for Positive (must have) and Negative (must not have) regexps.
// All checks are AND'ed.

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	provider "github.com/ContainerSolutions/argus/operator/internal/provider/schema"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Client struct {
	url                    string
	PositiveRegexp         string
	MinNumberPositiveMatch int
	NegativeRegesp         string
	MaxNumberNegativeMatch int
}

func newUnknownResult(reason string, err error) argusiov1alpha1.AttestationResult {
	r := argusiov1alpha1.AttestationResult{
		RunAt:  v1.Now(),
		Reason: reason,
		Result: argusiov1alpha1.AttestationResultTypeUnknown,
		Err:    err.Error(),
	}
	return r
}

func newFailResult(reason string) argusiov1alpha1.AttestationResult {
	r := argusiov1alpha1.AttestationResult{
		RunAt:  v1.Now(),
		Reason: reason,
		Result: argusiov1alpha1.AttestationResultTypeFail,
		Err:    "",
	}
	return r
}

func newPassResult(reason string) argusiov1alpha1.AttestationResult {
	r := argusiov1alpha1.AttestationResult{
		RunAt:  v1.Now(),
		Reason: reason,
		Result: argusiov1alpha1.AttestationResultTypePass,
		Err:    "",
	}
	return r
}
func (c *Client) Attest() (argusiov1alpha1.AttestationResult, error) {
	p := http.Client{}
	resp, err := p.Get(c.url)
	if err != nil {
		return newUnknownResult(fmt.Sprintf("could not GET url '%v'", c.url), err), err
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return newUnknownResult(fmt.Sprintf("could read response body for url '%v'", c.url), err), err
	}
	lines := bytes.Split(resBody, []byte("\n"))
	positive := regexp.MustCompile(c.PositiveRegexp)
	negative := regexp.MustCompile(c.NegativeRegesp)
	totalPos := 0
	totalNeg := 0
	reason := ""
	for i, line := range lines {
		if c.PositiveRegexp != "" {
			idx := positive.FindAllIndex(line, -1)
			totalPos = totalPos + len(idx)
			for _, t := range idx {
				reason = reason + fmt.Sprintf("Positive Match: line %v - %v\n", i+1, string(line[t[0]:t[1]]))
			}
		}
		if c.NegativeRegesp != "" {
			negIdx := negative.FindAllIndex(line, -1)
			totalNeg = totalNeg + len(negIdx)
			for _, t := range negIdx {
				reason = reason + fmt.Sprintf("Negative Match: line %v - %v\n", i+1, string(line[t[0]:t[1]]))
			}
		}
	}
	if totalPos >= c.MinNumberPositiveMatch && totalNeg <= c.MaxNumberNegativeMatch {
		return newFailResult(reason), nil
	}
	return newPassResult(reason), nil
}

func (c *Client) Close() error {
	return nil
}

type Provider struct{}

func (p *Provider) New(name string, spec *argusiov1alpha1.AttestationProviderSpec) (provider.AttestationClient, error) {
	var posMin string
	var negMax string
	url, ok := spec.ProviderConfig["url"]
	if !ok {
		return nil, fmt.Errorf("property 'url' is mandatory")
	}
	posRegexp, pOk := spec.ProviderConfig["positiveRegexp"]
	negRegexp, nOk := spec.ProviderConfig["negativeRegexp"]
	if !pOk && !nOk {
		return nil, fmt.Errorf("at least one of 'positiveRegexp' or 'negativeRegexp' is required")
	}
	posMin, ok = spec.ProviderConfig["minPositiveMatches"]
	if !ok {
		posMin = "1"
	}
	negMax, ok = spec.ProviderConfig["maxNegativeMatches"]
	if !ok {
		negMax = "0"
	}
	minNumberPositiveMatch, err := strconv.ParseInt(posMin, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("expected integer in 'minPositiveMatches': %w", err)
	}
	MaxNumberNegativeMatch, err := strconv.ParseInt(negMax, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("expected integer in 'maxNegativeMatches': %w", err)
	}
	c := &Client{
		url:                    url,
		PositiveRegexp:         posRegexp,
		NegativeRegesp:         negRegexp,
		MinNumberPositiveMatch: int(minNumberPositiveMatch),
		MaxNumberNegativeMatch: int(MaxNumberNegativeMatch),
	}
	return c, nil
}

func init() {
	provider.Register(&Provider{}, "file")
}
