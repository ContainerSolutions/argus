/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"argus/pkg/models"
	"argus/pkg/parser"
	"argus/pkg/resolver"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
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
		p, err := parser.Parse(c)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse config file '%v': %v\n", cfgFile, err)
			os.Exit(1)
		}
		d, err := json.Marshal(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not marshal parsed config file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(d))
		resolver.ResolveRequirements(p)
		d, _ = json.Marshal(p)
		fmt.Println(string(d))
		resolver.ResolveImplementations(p)
		d, _ = json.Marshal(p)
		fmt.Println(string(d))
		resolver.ResolveAttestations(p)
		d, _ = json.Marshal(p)
		fmt.Println(string(d))
		for _, r := range p.Resources {
			for _, a := range r.VerifiableAttestations {
				fmt.Printf("Verifying attestation %v\n", a.Name)
				cmd := exec.Command(a.CommandRef.Command, a.CommandRef.Args...)
				out, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Println("Failed!")
				}
				if cmd.ProcessState.ExitCode() != a.CommandRef.ExpectedExitCode {
					fmt.Printf("Code failed! Got %v But Expected %v\n", cmd.ProcessState.ExitCode(), a.CommandRef.ExpectedExitCode)
				}
				if a.CommandRef.ExpectedOutput != "" {
					if !strings.Contains(string(out), a.CommandRef.ExpectedOutput) {
						fmt.Printf("Output failed! Got %v\n", string(out))
					}
				}
				fmt.Println("verified!")
			}
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&cfgFile, "config", "c", ".argus-config.yaml", "Configuration file to run")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	// Use config file from the flag.
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.

}
