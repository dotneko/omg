/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"

	"github.com/spf13/cobra"
)

// balancesCmd represents the balances command
var balancesCmd = &cobra.Command{
	Use:   "balances [alias]",
	Short: "Query balances for an account",
	Long:  `Query balances for an account.`,
	Run: func(cmd *cobra.Command, args []string) {
		l := &omg.Accounts{}

		if err := l.Load(cfg.OmgFilename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		address := l.GetAddress(args[0])
		if address == "" {
			fmt.Printf("Error: account %q not found.\n", args[0])
			os.Exit(1)
		}
		balance, err := omg.GetBalanceAmount(address)
		if err != nil {
			fmt.Sprintln(err)
		}
		fmt.Printf("Balance = %.0f%s (~%.10f %s)\n", balance, cfg.Denom, omg.DenomToToken(balance), cfg.Token)
	},
}

func init() {
	queryCmd.AddCommand(balancesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// balancesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// balancesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
