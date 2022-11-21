/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io"
	"os"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"

	"github.com/spf13/cobra"
)

// balancesCmd represents the balances command
var balancesCmd = &cobra.Command{
	Aliases: []string{"bal", "b"},
	Use:     "balances [alias | address]",
	Short:   "Query balances for an account or address",
	Long:    `Query balances for an account or address.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		allAccounts, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}
		return balancesAction(os.Stdout, allAccounts, args)
	},
}

func init() {
	rootCmd.AddCommand(balancesCmd)

	balancesCmd.Flags().BoolP("all", "a", false, "Check all accounts in address book")
}

func balancesAction(out io.Writer, allAccounts bool, args []string) error {
	l := &omg.Accounts{}

	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	if allAccounts {
		for _, acc := range *l {
			if omg.IsNormalAddress(acc.Address) {
				balance, err := omg.GetBalanceDec(acc.Address)
				if err != nil {
					return err
				}
				fmt.Fprintf(out, "%10s [%s]: %s %s (%s %s)\n", acc.Alias, acc.Address, balance.String(), cfg.Denom, omg.DenomToTokenDec(balance).String(), cfg.Token)
			}
		}
		return nil
	}
	if len(args) == 0 {
		return fmt.Errorf("no account provided")
	}
	var address string
	var header string
	if omg.IsNormalAddress(args[0]) {
		address = args[0]
		header = ""
	} else {
		address = l.GetAddress(args[0])
		header = fmt.Sprintf("%10s [%s]: ", args[0], address)
	}
	if address == "" {
		return fmt.Errorf("account %q not found", args[0])
	}
	balance, err := omg.GetBalanceDec(address)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, header+"%s %s (%s %s)\n", balance.String(), cfg.Denom, omg.DenomToTokenDec(balance).String(), cfg.Token)
	return nil
}
