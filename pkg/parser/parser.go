package parser

import (
	"argus/pkg/models"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

func Parse(config models.ConfigFile) (*models.Configuration, error) {
	c := models.Configuration{}
	var err error
	c.Attestations, err = ParseAttestations(config.AttestationPath)
	if err != nil {
		return nil, err
	}
	c.Requirements, err = ParseRequirements(config.RequirementPath)
	if err != nil {
		return nil, err
	}
	c.Resources, err = ParseResources(config.ResourcePath)
	if err != nil {
		return nil, err
	}
	c.Implementations, err = ParseImplementations(config.ImplementationPath)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func ParseResources(config string) ([]models.Resource, error) {
	resources := []models.Resource{}
	walkFn := func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" {
			yamlFile, err := os.Open(path)
			if err != nil {
				return err
			}

			defer yamlFile.Close()

			byteValue, _ := io.ReadAll(yamlFile)

			var resource models.Resource
			if err := yaml.Unmarshal(byteValue, &resource); err != nil {
				return fmt.Errorf("Could not unmarshal file %s: %w", path, err)
			}
			resources = append(resources, resource)
		}
		return nil
	}
	err := filepath.Walk(config, walkFn)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func ParseRequirements(config string) ([]models.Requirement, error) {
	res := []models.Requirement{}
	walkFn := func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" {
			yamlFile, err := os.Open(path)
			if err != nil {
				return err
			}

			defer yamlFile.Close()

			byteValue, _ := io.ReadAll(yamlFile)

			var p models.Requirement
			if err := yaml.Unmarshal(byteValue, &p); err != nil {
				return fmt.Errorf("Could not unmarshal file %s: %w", path, err)
			}
			res = append(res, p)
		}
		return nil
	}
	err := filepath.Walk(config, walkFn)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ParseImplementations(config string) ([]models.Implementation, error) {
	res := []models.Implementation{}
	walkFn := func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" {
			yamlFile, err := os.Open(path)
			if err != nil {
				return err
			}

			defer yamlFile.Close()

			byteValue, _ := io.ReadAll(yamlFile)

			var p models.Implementation
			if err := yaml.Unmarshal(byteValue, &p); err != nil {
				return fmt.Errorf("Could not unmarshal file %s: %w", path, err)
			}
			res = append(res, p)
		}
		return nil
	}
	err := filepath.Walk(config, walkFn)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ParseAttestations(config string) ([]models.Attestation, error) {
	res := []models.Attestation{}
	walkFn := func(path string, info fs.FileInfo, err error) error {
		fmt.Println(path)
		if filepath.Ext(path) == ".yaml" {
			yamlFile, err := os.Open(path)
			if err != nil {
				return err
			}

			defer yamlFile.Close()

			byteValue, _ := io.ReadAll(yamlFile)

			var p models.Attestation
			if err := yaml.Unmarshal(byteValue, &p); err != nil {
				return fmt.Errorf("Could not unmarshal file %s: %w", path, err)
			}
			res = append(res, p)
		}
		return nil
	}
	err := filepath.Walk(config, walkFn)
	if err != nil {
		return nil, err
	}
	return res, nil
}
