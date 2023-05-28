/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"argus/pkg/parser"
	"argus/pkg/resolver"
	"argus/pkg/storage"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		resolver.Resolve(p)
		err = db.Save(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not save configuration '%v': %v\n", p, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
