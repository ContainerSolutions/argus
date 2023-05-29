/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"argus/pkg/results"
	"argus/pkg/storage"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var output string
var mode string

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
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
			fmt.Fprintf(os.Stderr, "could not configure database: %v\n", err)
			os.Exit(1)
		}
		switch mode {
		case "summary":
			err := results.Summary(config, output)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not generate report: %v\n", err)
				os.Exit(1)
			}
		case "detailed":
			err := results.Detailed(config, output)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not generate report: %v\n", err)
				os.Exit(1)
			}
		case "all":
			err := results.All(config, output)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not generate report: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, "'%v' is not a valid report type\n", mode)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringVarP(&mode, "mode", "m", "summary", "type of report. Possible values are 'summary' or 'detailed'")
	reportCmd.Flags().StringVarP(&output, "output", "o", "tsv", "command output. possible values are 'tsv' or 'json'")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
