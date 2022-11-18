/*
Copyright © 2022 dotneko

*/
package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"

	"github.com/spf13/cobra"
)

// delegateCmd represents the delegate command
var delegateCmd = &cobra.Command{
	Use:   "delegate [account] [validator] [amount][denom]",
	Short: "delegate [account] [validator] [amount][denom]",
	Long:  `Delegate tokens from account to validator.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		keyring, err := cmd.Flags().GetString("keyring")
		if err != nil {
			return err
		}
		auto, err := cmd.Flags().GetBool("auto")
		if err != nil {
			return err
		}
		return delegateAction(os.Stdout, keyring, auto, args)
	},
}

func init() {
	txCmd.AddCommand(delegateCmd)

	delegateCmd.Flags().StringP("remainder", "r", "1nom", "Remainder")
}

func delegateAction(out io.Writer, keyring string, auto bool, args []string) error {

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
	// Check balance for delegator
	fmt.Printf("Delegator %s [%s]\n", delegator, delegatorAddress)
	omg.CheckBalances(delegatorAddress)

	amount, err := omg.GetAmount(os.Stdin, "delegate", delegatorAddress, flag.Args()...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	omg.TxDelegateToValidator(delegator, valAddress, amount, keyring, auto)

	return nil
}
