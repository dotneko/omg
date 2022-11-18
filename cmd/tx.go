/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	cfg "github.com/dotneko/omg/config"
	"github.com/spf13/cobra"
)

// txCmd represents the tx command
var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "Execute a transaction",
	Long:  `Execute a transaction.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("tx called")
	// },
}

func init() {
	rootCmd.AddCommand(txCmd)

	txCmd.PersistentFlags().BoolP("auto", "a", false, "Auto confirm transaction")
	txCmd.PersistentFlags().StringP("keyring", "k", cfg.KeyringBackend, "Specify keyring-backend to use.")

}
