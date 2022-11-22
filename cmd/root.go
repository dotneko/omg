/*
Copyright © 2022 dotneko

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	version string = "0.3.1"
	rootCmd        = &cobra.Command{
		Use:          "omg",
		Short:        "Convenience wrapper for the Onomy Protocol",
		Long:         `omg - A CLI tool for interacting with the Onomy Protocol blockchain`,
		Version:      version,
		SilenceUsage: true,
	}
)

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

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	versionTemplate := `{{printf "%s - v%s\n" .Name .Version}}`
	rootCmd.SetVersionTemplate(versionTemplate)
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

func initConfig() {
}
