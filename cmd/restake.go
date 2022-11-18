/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"
	"github.com/spf13/cobra"
)

// restakeCmd represents the restake command
var restakeCmd = &cobra.Command{
	Use:   "restake [alias] [validator alias]",
	Short: "restake [alias] [validator alias]",
	Long:  `Withdraw all rewards for account, then re-delegate to validator.`,
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyring, err := cmd.Flags().GetString("keyring")
		if err != nil {
			return err
		}
		auto, err := cmd.Flags().GetBool("auto")
		if err != nil {
			return err
		}
		remainder, err := cmd.Flags().GetString("remainder")
		if err != nil {
			return err
		}
		return restakeAction(os.Stdout, remainder, keyring, auto, args)
	},
}

func init() {
	txCmd.AddCommand(restakeCmd)

	restakeCmd.Flags().StringP("remainder", "r", "1nom", "Remainder after restake")
}

func restakeAction(out io.Writer, remainder string, keyring string, auto bool, args []string) error {
	// Ensure all arguments provided
	if len(args) != 2 {
		return fmt.Errorf("Expecting [delegator] [validator]")
	}
	delegator := args[0]
	validator := args[1]

	remainAmt, denom, err := omg.StrSplitAmountDenom(remainder)
	if err != nil {
		return err
	}
	if denom == cfg.Token {
		remainAmt = omg.TokenToDenom(remainAmt)
		denom = cfg.Denom
	} else if denom != cfg.Denom {
		return fmt.Errorf("Denomination not specified")
	}
	l := &omg.Accounts{}
	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	delegatorAddress := l.GetAddress(delegator)

	if !omg.IsNormalAddress(delegatorAddress) {
		return fmt.Errorf("Invalid delegator account: %s\n", delegatorAddress)
	}
	var valAddress string
	if omg.IsValidatorAddress(validator) {
		valAddress = validator
	} else {
		valAddress = l.GetAddress(validator)

		if !omg.IsValidatorAddress(valAddress) {
			return fmt.Errorf("Invalid validator: %q \n", valAddress)
		}
	}

	// Check balance for delegator
	fmt.Fprintf(out, "Delegator %s [%s]\n", delegator, delegatorAddress)
	balanceBefore, err := omg.GetBalanceAmount(delegatorAddress)
	if err != nil {
		return fmt.Errorf("Error getting balance for %s\n", delegator)
	}
	fmt.Fprintf(out, "Existing balance: %.0f %s\n", balanceBefore, cfg.Denom)
	fmt.Fprintf(out, "Withdrawing rewards for %s [%s]\n", delegator, delegatorAddress)
	omg.TxWithdrawRewards(out, delegator, keyring, auto)

	// Wait till balance is updated
	var balance *float64
	balance = new(float64)
	count := 0
	for count <= 10 {
		*balance, _ = omg.GetBalanceAmount(delegatorAddress)
		if *balance > balanceBefore {
			fmt.Fprintf(out, "Updated balance  : %.0f %s\n", *balance, cfg.Denom)
			break
		}
		count++
		time.Sleep(1 * time.Second)
	}
	// If balance not updated and -auto flag set then abort
	if auto && *balance == balanceBefore {
		return fmt.Errorf("Balance not increased. Aborting auto-restake")
	}
	// Restake amount leaving approx remainder of 1 token
	amount := *balance - remainAmt
	omg.TxDelegateToValidator(delegator, valAddress, amount, keyring, auto)
	return nil
}
