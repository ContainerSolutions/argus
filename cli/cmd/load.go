/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/ContainerSolutions/argus/cli/pkg/parser"
	"github.com/ContainerSolutions/argus/cli/pkg/resolver"
	"github.com/ContainerSolutions/argus/cli/pkg/storage"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Loads configuration and saves configuration to state database. Run this first.",
	Long: `Loads the configuration file and saves this configuration to the configured file.
Must be run before program will work.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := loadConfig()
		p, err := parser.Parse(c)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse config file '%v': %v\n", cfgFile, err)
			os.Exit(1)
		}
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
		_, err = resolver.Resolve(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not resolve configuration '%v': %v\n", p, err)
			os.Exit(1)
		}
		err = db.Save(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not save configuration '%v': %v\n", p, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
}
