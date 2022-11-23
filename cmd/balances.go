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
		detail, err := cmd.Flags().GetBool("detail")
		if err != nil {
			return err
		}
		raw, err := cmd.Flags().GetBool("raw")
		if err != nil {
			return err
		}
		return balancesAction(os.Stdout, allAccounts, detail, raw, args)
	},
}

func init() {
	rootCmd.AddCommand(balancesCmd)

	balancesCmd.Flags().BoolP("all", "a", false, "Check all accounts in address book")
	balancesCmd.Flags().BoolP("detail", "d", false, "Detailed output")
	balancesCmd.Flags().BoolP("raw", "r", false, "Raw output")

}

func balancesAction(out io.Writer, allAccounts bool, detail bool, raw bool, args []string) error {
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
				if raw {
					fmt.Fprintf(out, "%10s %s%s\n", acc.Address, balance.String(), cfg.BaseDenom)
				} else if detail {
					fmt.Fprintf(out, "%10s [%s]: %30s %s (%s %s)\n", acc.Alias, acc.Address, omg.PrettifyDenom(balance), cfg.BaseDenom, omg.DenomToTokenDec(balance).String(), cfg.Token)
				} else {
					fmt.Fprintf(out, "%10s : %30s %s (%s %s)\n", acc.Alias, omg.PrettifyDenom(balance), cfg.BaseDenom, omg.DenomToTokenDec(balance).String(), cfg.Token)
				}
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
		if detail {
			header = fmt.Sprintf("%10s [%s]: ", args[0], address)
		} else {
			header = fmt.Sprintf("%10s : ", args[0])
		}
	}
	if address == "" {
		return fmt.Errorf("account %q not found", args[0])
	}
	balance, err := omg.GetBalanceDec(address)
	if err != nil {
		return err
	}
	if raw {
		fmt.Fprintf(out, "%s%s", balance.String(), cfg.BaseDenom)
	} else {
		fmt.Fprintf(out, header+"%s %s (%s %s)\n", omg.PrettifyDenom(balance), cfg.BaseDenom, omg.DenomToTokenDec(balance).String(), cfg.Token)
	}
	return nil
}
