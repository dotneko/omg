/*
Copyright © 2022 dotneko

*/
package cmd

import (
	"fmt"
	"io"
	"os"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"
	"github.com/shopspring/decimal"

	"github.com/spf13/cobra"
)

// delegateCmd represents the delegate command
var delegateCmd = &cobra.Command{
	Aliases: []string{"del", "d"},
	Use:     "delegate [account] [validator] [amount][denom]",
	Short:   "Delegate tokens from account to validator",
	Long:    `Delegate tokens from account to validator.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		keyring, err := cmd.Flags().GetString("keyring")
		if err != nil {
			return err
		}
		auto, err := cmd.Flags().GetBool("auto")
		if err != nil {
			return err
		}
		all, err := cmd.Flags().GetBool("full")
		if err != nil {
			return err
		}
		remainder, err := cmd.Flags().GetString("remainder")
		if err != nil {
			return err
		}
		return delegateAction(os.Stdout, keyring, auto, all, remainder, args)
	},
}

func init() {
	txCmd.AddCommand(delegateCmd)

	delegateCmd.Flags().BoolP("full", "f", false, "Delegate full balance amount (less remainder)")
	delegateCmd.Flags().StringP("remainder", "r", "100000000anom", "Remainder")
}

func delegateAction(out io.Writer, keyring string, auto bool, all bool, remainder string, args []string) error {

	delegator, validator, err := omg.GetTxAccounts(os.Stdin, "delegate", args...)
	if err != nil {
		return err
	}

	l := &omg.Accounts{}
	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}

	// Check if delegator in list and is not validator account
	delegatorAddress := l.GetAddress(delegator)

	if !omg.IsNormalAddress(delegatorAddress) {
		return fmt.Errorf("invalid delegator address: %s", delegatorAddress)
	}
	// Check if valid validator address
	valAddress := l.GetAddress(validator)

	if !omg.IsValidatorAddress(valAddress) {
		return fmt.Errorf("invalid validator address %s", valAddress)
	}
	fmt.Fprintf(out, "Delegator         : %s [%s]\n", delegator, delegatorAddress)
	balance, err := omg.GetBalanceDec(delegatorAddress)
	if err != nil {
		return err
	}
	var (
		amount    decimal.Decimal
		remainAmt decimal.Decimal
		denom     string
	)
	if all {
		// Parse remainder
		remainAmt, denom, err = omg.StrSplitAmountDenomDec(remainder)
		if err != nil {
			return err
		}
		if denom == cfg.Token {
			amount = omg.TokenToDenomDec(remainAmt)
			denom = cfg.Denom
		} else if denom != cfg.Denom {
			return fmt.Errorf("denomination not specified")
		} else {
			amount = balance.Sub(remainAmt)
		}
		if amount.LessThanOrEqual(decimal.NewFromInt(0)) {
			return fmt.Errorf("insufficient balance after deducting remainder: %s %s", amount.String(), denom)
		}
	} else {
		amount, err = omg.GetAmount(os.Stdin, "delegate", delegatorAddress, args...)
		if err != nil {
			return err
		}
		remainAmt = balance.Sub(amount)
	}
	fmt.Fprintf(out, "Available balance : %s%s\n", omg.PrettifyDenom(balance), denom)
	fmt.Fprintf(out, "Delegation amount : %s%s\n", omg.PrettifyDenom(amount), denom)
	fmt.Fprintf(out, "Remainder amount  : %s%s\n", omg.PrettifyDenom(remainAmt), denom)
	fmt.Fprintln(out, "----")
	omg.TxDelegateToValidator(delegator, valAddress, amount, keyring, auto)

	return nil
}
