/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/ContainerSolutions/argus/cli/pkg/models"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "argus",
	Short: "Cloud Control Framework CLI",
	Long:  `Cloud Control Framework CLI`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.argus.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", ".argus-config.yaml", "Configuration file to run")
}

func initConfig() {
	// Use config file from the flag.
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.

}

func loadConfig() *models.ConfigFile {
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
	return &c
}
