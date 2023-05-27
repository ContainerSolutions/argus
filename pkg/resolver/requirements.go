package resolver

import "argus/pkg/models"

func ResolveRequirements(config *models.Configuration) (*models.Configuration, error) {
	for k, r := range config.Resources {
		resolveReqForResource(&r, config.Requirements)
		config.Resources[k] = r
	}
	return config, nil
}

func ResolveImplementations(config *models.Configuration) (*models.Configuration, error) {
	for k, r := range config.Resources {
		resolveImpForResource(&r, config.Implementations)
		config.Resources[k] = r
	}
	return config, nil

}

func ResolveAttestations(config *models.Configuration) (*models.Configuration, error) {
	for k, r := range config.Resources {
		resolveAttForResource(&r, config.Attestations)
		config.Resources[k] = r
	}
	return config, nil
}

func resolveAttForResource(current *models.Resource, attestations []models.Attestation) {
	for _, attestation := range attestations {
		for _, implementation := range current.ApplicableImplementations {
			if implementation.AttestationRef == attestation.Name {
				current.VerifiableAttestations = append(current.VerifiableAttestations, &attestation)
			}
		}
	}
}
func resolveImpForResource(current *models.Resource, implementations []models.Implementation) {
	for _, implementation := range implementations {
		for _, requirement := range current.ApplicableRequirements {
			if implementation.RequirementRef.Code == requirement.Code && implementation.RequirementRef.Version == requirement.Version && (implementation.ResourceRef == current.Name || contains(current.Parents, implementation.ResourceRef)) {
				current.ApplicableImplementations = append(current.ApplicableImplementations, &implementation)
			}
		}
	}
}
func resolveReqForResource(current *models.Resource, requirements []models.Requirement) {
	for _, requirement := range requirements {
		for _, class := range current.Classes {
			if contains(requirement.ApplicableResourceClasses, class) {
				current.ApplicableRequirements = append(current.ApplicableRequirements, &requirement)
			}
		}
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
