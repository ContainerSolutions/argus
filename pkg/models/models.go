package models

// add json labels to the following structures
type Resource struct {
	Name                      string   `json:"name"`
	Type                      string   `json:"type"`
	Classes                   []string `json:"classes"`
	Parents                   []string `json:"parents"`
	ApplicableRequirements    []*Requirement
	ApplicableImplementations []*Implementation
	VerifiableAttestations    []*Attestation
}

type Requirement struct {
	Name                          string   `json:"name"`
	Version                       string   `json:"version"`
	Code                          string   `json:"code"`
	Class                         string   `json:"class"`
	Category                      string   `json:"category"`
	ApplicableResourceClasses     []string `json:"applicableResourceClasses"`
	RequiredImplementationClasses []string `json:"requiredImplementationClasses"`
	requirementVersion            string
}

type RequirementRef struct {
	Code    string `json:"code"`
	Version string `json:"version"`
}
type Implementation struct {
	Name                    string         `json:"name"`
	Class                   string         `json:"class"`
	RequirementRef          RequirementRef `json:"requirementRef"`
	ResourceRef             string         `json:"resourceRef"`
	boundRequirementVersion string
	AttestationRef          string `json:"attestationRef"`
}

type Attestation struct {
	Name       string               `json:"name"`
	Type       string               `json:"type"`
	CommandRef AttestationByCommand `json:"commandRef"`
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
	ResourcePath       string `json:"resourcePath"`
	RequirementPath    string `json:"requirementPath"`
	ImplementationPath string `json:"implementationPath"`
	AttestationPath    string `json:"attestationPath"`
}
