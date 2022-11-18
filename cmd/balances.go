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
	Use:     "balances [alias]",
	Short:   "balances [alias] or 'balances -a' for all accounts",
	Long:    `Query balances for an account or all normal accounts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		allAccounts, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}
		return balancesAction(os.Stdout, allAccounts, args)
	},
}

func init() {
	queryCmd.AddCommand(balancesCmd)

	balancesCmd.Flags().BoolP("all", "a", false, "List balances for all accounts")
}

func balancesAction(out io.Writer, allAccounts bool, args []string) error {
	l := &omg.Accounts{}

	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	if allAccounts {
		for _, acc := range *l {
			if omg.IsNormalAddress(acc.Address) {
				balance, err := omg.GetBalanceAmount(acc.Address)
				if err != nil {
					return err
				}
				fmt.Printf("%10s [%s]: %.0f%s (~%.10f %s)\n", acc.Alias, acc.Address, balance, cfg.Denom, omg.DenomToToken(balance), cfg.Token)
			}
		}
		return nil
	}
	if len(args) == 0 {
		return fmt.Errorf("No account provided.\n")
	}
	address := l.GetAddress(args[0])
	if address == "" {
		return fmt.Errorf("Account %q not found.\n", args[0])
	}
	balance, err := omg.GetBalanceAmount(address)
	if err != nil {
		return err
	}
	fmt.Printf("%10s [%s]: %.0f%s (~%.10f %s)\n", args[0], address, balance, cfg.Denom, omg.DenomToToken(balance), cfg.Token)
	return nil
}
