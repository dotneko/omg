/*
Copyright © 2022 dotneko
*/
package cmd

import (
	cfg "github.com/dotneko/omg/config"
	"github.com/spf13/cobra"
)

// txCmd represents the tx command
var txCmd = &cobra.Command{
	Aliases: []string{"t"},
	Use:     "tx",
	Short:   "Execute a transaction",
	Long:    `Execute a transaction.`,
}

func init() {
	rootCmd.AddCommand(txCmd)

	txCmd.PersistentFlags().BoolP("yes", "y", false, "Auto confirm transaction")
	txCmd.PersistentFlags().StringP("keyring", "k", cfg.KeyringBackend, "Specify keyring-backend to use")
	txCmd.PersistentFlags().BoolP("txhash", "t", false, "Output transaction hash only when auto confirm")

}
