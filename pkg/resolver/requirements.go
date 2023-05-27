package resolver

import (
	"argus/pkg/models"
	"argus/pkg/utils"
)

func Resolve(config *models.Configuration) (*models.Configuration, error) {
	var err error
	config, err = ResolveRequirements(config)
	if err != nil {
		return nil, err
	}
	config, err = ResolveImplementations(config)
	if err != nil {
		return nil, err
	}
	return ResolveAttestations(config)
}
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
		for k, reqBlock := range current.Requirements {
			for kk, impBlock := range reqBlock.Implementations {
				implementation := impBlock.Implementaiton
				if implementation.AttestationRef == attestation.Name {
					impBlock.Attestation = &attestation
				}
				reqBlock.Implementations[kk] = impBlock
			}
			current.Requirements[k] = reqBlock
		}
	}
}
func resolveImpForResource(current *models.Resource, implementations []models.Implementation) {
	for _, implementation := range implementations {
		for k, reqBlock := range current.Requirements {
			requirement := reqBlock.Requirement
			if implementation.RequirementRef.Code == requirement.Code && implementation.RequirementRef.Version == requirement.Version && (implementation.ResourceRef == current.Name || utils.Contains(current.Parents, implementation.ResourceRef)) {
				if reqBlock.Implementations == nil {
					reqBlock.Implementations = make(map[string]models.ImplementationBlock)
				}
				reqBlock.Implementations[implementation.Name] = models.ImplementationBlock{
					Implementaiton: &implementation,
				}
			}
			current.Requirements[k] = reqBlock
		}
	}
}
func resolveReqForResource(current *models.Resource, requirements []models.Requirement) {
	for _, requirement := range requirements {
		for _, class := range current.Classes {
			if utils.Contains(requirement.ApplicableResourceClasses, class) {
				if current.Requirements == nil {
					current.Requirements = make(map[string]models.RequirementBlock)
				}
				current.Requirements[requirement.Name] = models.RequirementBlock{
					Requirement: &requirement,
				}
			}
		}
	}
}
