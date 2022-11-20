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

	"github.com/spf13/cobra"
)

// delegateCmd represents the delegate command
var delegateCmd = &cobra.Command{
	Aliases: []string{"del"},
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
		return fmt.Errorf("Invalid delegator address: %s\n", delegatorAddress)
	}
	// Check if valid validator address
	valAddress := l.GetAddress(validator)

	if !omg.IsValidatorAddress(valAddress) {
		return fmt.Errorf("Invalid validator address %s\n", valAddress)
	}
	fmt.Fprintf(out, "Delegator         : %s [%s]\n", delegator, delegatorAddress)
	balance, err := omg.GetBalanceAmount(delegatorAddress)
	if err != nil {
		return err
	}
	var (
		amount    float64
		remainAmt float64
		denom     string
	)
	if all {
		// Parse remainder
		remainAmt, denom, err = omg.StrSplitAmountDenom(remainder)
		if err != nil {
			return err
		}
		if denom == cfg.Token {
			amount = omg.TokenToDenom(remainAmt)
			denom = cfg.Denom
		} else if denom != cfg.Denom {
			return fmt.Errorf("Denomination not specified")
		}
		amount = balance - remainAmt
		if amount <= 0 {
			return fmt.Errorf("Insufficient balance after deducting remainder: %.0f", amount)
		}
	} else {
		amount, err = omg.GetAmount(os.Stdin, "delegate", delegatorAddress, args...)
		if err != nil {
			return err
		}
		remainAmt = balance - amount
	}
	fmt.Fprintf(out, "Available balance : %s\n", omg.PrettifyDenom(balance))
	fmt.Fprintf(out, "Delegation amount : %s\n", omg.PrettifyDenom(amount))
	fmt.Fprintf(out, "Remander amount   : %s\n", omg.PrettifyDenom(remainAmt))
	fmt.Fprintln(out, "----")
	omg.TxDelegateToValidator(delegator, valAddress, amount, keyring, auto)

	return nil
}
