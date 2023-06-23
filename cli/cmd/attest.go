/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/ContainerSolutions/argus/cli/pkg/attester"
	"github.com/ContainerSolutions/argus/cli/pkg/results"
	"github.com/ContainerSolutions/argus/cli/pkg/storage"
	"github.com/ContainerSolutions/argus/cli/pkg/utils"

	"github.com/spf13/cobra"
)

// attestCmd represents the attest command
var attestCmd = &cobra.Command{
	Use:   "attest",
	Short: "Run an attestation",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := loadConfig()
		db, err := storage.Init(c.Driver)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not initialize database: %v\n", err)
			os.Exit(1)
		}
		err = db.Configure(c.DriverConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not configure database: %v\n", err)
			os.Exit(1)
		}
		config, err := db.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not load database: %v\n", err)
			os.Exit(1)
		}
		for kkk, r := range config.Resources {
			implementedRequirements := 0
			for kk, req := range r.Requirements {
				totalImplementations := 0
				attestedImplementations := 0
				for k, i := range req.Implementations {
					verifiedAttestations := 0
					if utils.Contains(req.Requirement.RequiredImplementationClasses, i.Implementation.Class) {
						totalImplementations = totalImplementations + 1
					}
					for _, a := range i.Attestation {
						fmt.Printf("Resource: %v\nRequirement:'%v'\nImplementation: '%v'\nAttestation: '%v'\nResult: ", r.Name, req.Requirement.Name, i.Implementation.Name, a.Attestation.Name)
						attester, _ := attester.Init(a.Attestation.Type)
						res, err := attester.Attest(a.Attestation)
						if err != nil {
							os.Exit(1)
						}
						a.Attested = res.Result == "PASS"
						fmt.Printf("%v\n\n", res.Result)
						if a.Attested {
							verifiedAttestations = verifiedAttestations + 1
						}
					}
					i.TotalAttestations = len(i.Attestation)
					i.Attested = false
					i.VerifiedAttestations = verifiedAttestations
					if i.VerifiedAttestations == i.TotalAttestations {
						i.Attested = true
					}
					if i.Attested && utils.Contains(req.Requirement.RequiredImplementationClasses, i.Implementation.Class) {
						attestedImplementations = attestedImplementations + 1
					}
					req.Implementations[k] = i
				}
				req.AttestedImplementations = attestedImplementations
				req.TotalImplementations = totalImplementations
				if len(req.Requirement.RequiredImplementationClasses) <= req.AttestedImplementations {
					req.Implemented = true
					implementedRequirements = implementedRequirements + 1
				}
				r.Requirements[kk] = req
			}
			r.ImplementedRequirements = implementedRequirements
			if len(r.Requirements) == implementedRequirements {
				r.Implemented = true
			}
			config.Resources[kkk] = r
		}
		err = db.Save(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not save db after attestation: %v\n", err)
			os.Exit(1)
		}
		err = results.Summary(config, "tsv")
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not generate summary: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(attestCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}