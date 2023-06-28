package models

import "time"

// add json labels to the following structures
type Resource struct {
	Name                    string   `json:"name"`
	Type                    string   `json:"type"`
	Classes                 []string `json:"classes"`
	Parents                 []string `json:"parents"`
	Requirements            map[string]RequirementBlock
	ImplementedRequirements int
	Implemented             bool
}

type RequirementBlock struct {
	Requirement             *Requirement
	Implementations         map[string]ImplementationBlock
	Implemented             bool
	AttestedImplementations int
	TotalImplementations    int
	RunAt                   string
}

type ImplementationBlock struct {
	Implementation       *Implementation
	Attestation          map[string]AttestationBlock
	TotalAttestations    int
	Attested             bool
	VerifiedAttestations int
	RunAt                string
}

type AttestationBlock struct {
	Attestation *Attestation
	Attested    bool
	RunAt       string
	Logs        string
}

type Requirement struct {
	Name                          string   `json:"name"`
	Version                       string   `json:"version"`
	Code                          string   `json:"code"`
	Class                         string   `json:"class"`
	Category                      string   `json:"category"`
	ApplicableResourceClasses     []string `json:"applicableResourceClasses"`
	RequiredImplementationClasses []string `json:"requiredImplementationClasses"`
}

type RequirementRef struct {
	Code    string `json:"code"`
	Version string `json:"version"`
}
type Implementation struct {
	Name           string         `json:"name"`
	Class          string         `json:"class"`
	RequirementRef RequirementRef `json:"requirementRef"`
	ResourceRef    []string       `json:"resourceRef"`
}

type Attestation struct {
	Name              string               `json:"name"`
	Type              string               `json:"type"`
	Result            AttestationResult    `json:"result"`
	CommandRef        AttestationByCommand `json:"commandRef"`
	ImplementationRef string               `json:"implementationRef"`
}

type AttestationByCommand struct {
	Command          string   `json:"command"`
	Args             []string `json:"args"`
	ExpectedExitCode int      `json:"expectedExitCode"`
	ExpectedOutput   string   `json:"expectedOutput,omitempty"`
}

type Configuration struct {
	Resources       []Resource       `json:"resources"`
	Requirements    []Requirement    `json:"requirements"`
	Implementations []Implementation `json:"implementations"`
	Attestations    []Attestation    `json:"attestations"`
}

type ConfigFile struct {
	ResourcePath       string                 `json:"resourcePath"`
	RequirementPath    string                 `json:"requirementPath"`
	ImplementationPath string                 `json:"implementationPath"`
	AttestationPath    string                 `json:"attestationPath"`
	Driver             string                 `json:"driver"`
	DriverConfig       map[string]interface{} `json:"driverConfig"`
}

type AttestationResult struct {
	Command string
	Logs    string
	Result  string
	Reason  string
	Err     string
	RunAt   time.Time
}
