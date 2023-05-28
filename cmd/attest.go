/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"argus/pkg/attester"
	"argus/pkg/results"
	"argus/pkg/storage"
	"argus/pkg/utils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// attestCmd represents the attest command
var attestCmd = &cobra.Command{
	Use:   "attest",
	Short: "A brief description of your command",
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
					if utils.Contains(req.Requirement.RequiredImplementationClasses, i.Implementaiton.Class) {
						totalImplementations = totalImplementations + 1
					}
					for _, a := range i.Attestation {
						attester, _ := attester.Init(a.Attestation.Type)
						res, err := attester.Attest(a.Attestation)
						if err != nil {
							os.Exit(1)
						}
						a.Attested = res.Result == "PASS"
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
					if i.Attested && utils.Contains(req.Requirement.RequiredImplementationClasses, i.Implementaiton.Class) {
						attestedImplementations = attestedImplementations + 1
					}
					req.Implementations[k] = i
				}
				req.AttestedImplementations = attestedImplementations
				req.TotalImplementaitons = totalImplementations
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
		err = results.Detailed(config, "tsv")
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
