/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

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
		w := tabwriter.NewWriter(os.Stdout, 10, 4, 2, ' ', 0)
		var line string
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
						line = fmt.Sprintf("Resource:\t%v\nRequirement:\t'%v'\nImplementation:\t'%v'\nAttestation:\t'%v'\nResult:\t", r.Name, req.Requirement.Name, i.Implementation.Name, a.Attestation.Name)
						_, err := w.Write([]byte(line))
						if err != nil {
							fmt.Fprintf(os.Stderr, "error happened while printing to output:%v", err)
						}
						attester, _ := attester.Init(a.Attestation.Type)
						res, err := attester.Attest(a.Attestation)
						if err != nil {
							os.Exit(1)
						}
						a.Attested = res.Result == "PASS"
						line = fmt.Sprintf("%v\n\n", res.Result)
						_, err = w.Write([]byte(line))
						if err != nil {
							fmt.Fprintf(os.Stderr, "error happened while printing to output:%v", err)
						}
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
		w.Flush()
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
}
