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
	Use:     "balances [name|address]",
	Short:   "Query balances for an account or address",
	Long:    `Query balances for an account or address.`,
	Args: func(cmd *cobra.Command, args []string) error {
		allAccounts, _ := cmd.Flags().GetBool("all")
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		if len(args) != 1 && !allAccounts {
			fmt.Printf("Error: expecting [name|address] or --all flag\n")
			cmd.Help()
			os.Exit(0)
		}
		return nil
	},
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
					fmt.Printf("Error: %s\n", err.Error())
					continue
				}
				fmt.Fprintf(out, "%10s [%s]: %s %s (%s %s)\n", acc.Alias, acc.Address, balance.String(), cfg.BaseDenom, omg.DenomToTokenDec(balance).String(), cfg.Token)
			}
		}
		return nil
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
	fmt.Fprintf(out, header+"%s %s (%s %s)\n", balance.String(), cfg.BaseDenom, omg.DenomToTokenDec(balance).String(), cfg.Token)
	return nil
}
