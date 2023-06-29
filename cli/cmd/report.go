/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/ContainerSolutions/argus/cli/pkg/results"
	"github.com/ContainerSolutions/argus/cli/pkg/storage"

	"github.com/spf13/cobra"
)

var output string
var mode string

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Run a report on the current status",
	Long: `Run a report on the current status. These are the modes available:

- summary   - a summary?
- detailed  - more detail?
- all       - ??
`,
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
}
