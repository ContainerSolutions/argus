/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"argus/pkg/attester"
	"argus/pkg/models"
	"argus/pkg/storage"
	"argus/pkg/utils"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		fmt.Println("attest called")
		fmt.Println("run called")
		viper.SetConfigFile(cfgFile)
		viper.AutomaticEnv() // read in environment variables that match
		if err := viper.ReadInConfig(); cfgFile != "" && err != nil {
			fmt.Fprintf(os.Stderr, "Could not read file '%v': %v\n", cfgFile, err)
			os.Exit(1)
		}
		c := models.ConfigFile{}
		err := viper.Unmarshal(&c)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not unmarshal config file '%v': %v\n", cfgFile, err)
			os.Exit(1)
		}
		fmt.Println(c)
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
		fmt.Println("Running Program")
		for kkk, r := range config.Resources {
			fmt.Println(r.Name)
			implementedRequirements := 0
			for kk, req := range r.Requirements {
				fmt.Println(req.Requirement.Name)
				fmt.Println(req.Implementations)
				totalImplementations := 0
				attestedImplementations := 0
				for k, i := range req.Implementations {
					if utils.Contains(req.Requirement.RequiredImplementationClasses, i.Implementaiton.Class) {
						totalImplementations = totalImplementations + 1

					}
					fmt.Println(i.Implementaiton.Name)
					a := i.Attestation
					fmt.Printf("Verifying attestation %v\n", a.Name)
					attester, _ := attester.Init(a.Type)
					res, err := attester.Attest(a)
					if err != nil {
						os.Exit(1)
					}
					i.Attested = res.Result == "PASS"
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
		fmt.Println("Attestation finished. Results:")
		d, _ := json.MarshalIndent(config.Resources, "", " ")
		fmt.Println(string(d))
		// fmt.Println()
		// fmt.Println(results.Summary(config))
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(attestCmd)
	attestCmd.Flags().StringVarP(&cfgFile, "config", "c", ".argus-config.yaml", "Configuration file to run")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
