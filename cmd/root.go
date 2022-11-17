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
	version string = "0.3.0"
	rootCmd        = &cobra.Command{
		Use:     "omg",
		Short:   "Convenience wrapper for the Onomy Protocol",
		Long:    `omg - A CLI tool for interacting with the Onomy Protocol blockchain`,
		Version: version,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
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
	// err := omg.ParseConfig("../")
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// }

	// viper.SetConfigName(".omgconfig.yaml")
	// viper.SetConfigType("yaml")
	// viper.AddConfigPath("./")

	// // Check given pathstr if valid path
	// // _, err := os.Stat(pathstr)
	// // if err == nil {
	// // 	viper.AddConfigPath(pathstr)
	// // }
	// // Check home directory
	// home, _ := os.UserHomeDir()
	// viper.AddConfigPath(home)

	// err := viper.ReadInConfig()
	// if err != nil {
	// 	fmt.Printf("Error reading configuration, %s", err.Error())
	// }
	// err = viper.Unmarshal()
	// if err != nil {
	// 	return fmt.Errorf(err.Error())
	// }
}
