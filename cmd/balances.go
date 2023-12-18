/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"

	"github.com/spf13/cobra"
)

// balancesCmd represents the balances command
var balancesCmd = &cobra.Command{
	Aliases: []string{"balance", "bal", "b"},
	Use:     "balances [name|address]",
	Short:   "Query balances for an account or address",
	Long:    `Query balances for an account or address.`,
	Args: func(cmd *cobra.Command, args []string) error {
		allAccounts, _ := cmd.Flags().GetBool("all")
		if len(args) < 1 && !allAccounts {
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
		less, err := cmd.Flags().GetString("less")
		if err != nil {
			return err
		}
		raw, err := cmd.Flags().GetBool("raw")
		if err != nil {
			return err
		}
		token, err := cmd.Flags().GetBool("token")
		if err != nil {
			return err
		}
		var outType string = ""
		if raw {
			outType = omg.RAW
		} else if token {
			outType = omg.TOKEN
		} else if detail {
			outType = omg.DETAIL
		}
		return balancesAction(os.Stdout, allAccounts, less, outType, args)
	},
}

func init() {
	rootCmd.AddCommand(balancesCmd)

	balancesCmd.Flags().BoolP("all", "a", false, "Check all accounts in address book")
	balancesCmd.Flags().BoolP("detail", "d", false, "Detailed output")
	balancesCmd.Flags().StringP("less", "l", "", "Output balance less [amount] for individual account")
	balancesCmd.Flags().BoolP("raw", "r", false, "Raw output")
	balancesCmd.Flags().BoolP("token", "t", false, "Token amount output")

}

func balancesAction(out io.Writer, allAccounts bool, less, outType string, args []string) error {

	l := &omg.Accounts{}
	if err := l.Load(cfg.OmgFilepath); err != nil {
		return err
	}

	if allAccounts {
		for _, acc := range *l {
			if omg.IsNormalAddress(acc.Address) {
				balance, err := omg.GetBalance(acc.Address)
				if err != nil {
					fmt.Fprintf(out, "Error: %20s : %s\n", acc.Alias, err.Error())
					continue
				}
				if outType == omg.RAW || outType == omg.TOKEN {
					fmt.Fprintf(out, "%s ", acc.Address)
				}
				omg.OutputAmount(out, acc.Alias, acc.Address, balance.String(), outType)
			}
		}
		return nil
	}

	var name, address string
	if omg.IsNormalAddress(args[0]) {
		name = ""
		address = args[0]
	} else {
		name = args[0]
		address = l.GetAddress(args[0])
	}
	if address == "" {
		return fmt.Errorf("account %q not found", args[0])
	}
	balance, err := omg.GetBalance(address)
	if err != nil {
		return err
	}
	if less != "" {
		normalizedLess, err := omg.NormalizeAmountDenom((less))
		if err != nil {
			return fmt.Errorf("error normalizing less amount %s: %s", less, err)
		}
		lessCoin, err := sdktypes.ParseCoinNormalized(normalizedLess)
		if err != nil {
			return fmt.Errorf("error parsing less amount %s: %s", less, err)
		}
		finalBalance := balance.Sub(lessCoin)
		if outType != omg.RAW && outType != omg.TOKEN {
			fmt.Fprintf(out, "Balance (%s %s) less (%s):\n", omg.PrettifyBaseAmt(balance.String()), cfg.BaseDenom, omg.PrettifyBaseAmt(lessCoin.String()))
		}
		omg.OutputAmount(out, name, address, finalBalance.String(), outType)
	} else {
		omg.OutputAmount(out, name, address, balance.String(), outType)
	}
	return nil
}
