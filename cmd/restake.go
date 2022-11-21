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
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

// restakeCmd represents the restake command
var restakeCmd = &cobra.Command{
	Use:   "restake [name] [validator alias]",
	Short: "Restake rewards for account to validator",
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

	restakeCmd.Flags().StringP("remainder", "r", "100000000000000anom", "Remainder after restake")
}

func restakeAction(out io.Writer, remainder string, keyring string, auto bool, args []string) error {
	// Ensure all arguments provided
	if len(args) != 2 {
		return fmt.Errorf("expecting [delegator] [validator]")
	}
	delegator := args[0]
	validator := args[1]

	remainAmt, denom, err := omg.StrSplitAmountDenomDec(remainder)
	if err != nil {
		return err
	}
	if denom == cfg.Token {
		remainAmt = omg.TokenToDenomDec(remainAmt)
		denom = cfg.Denom
	} else if denom != cfg.Denom {
		return fmt.Errorf("denomination not specified")
	}
	l := &omg.Accounts{}
	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	delegatorAddress := l.GetAddress(delegator)

	if !omg.IsNormalAddress(delegatorAddress) {
		return fmt.Errorf("invalid delegator account: %s", delegatorAddress)
	}
	var valAddress string
	if omg.IsValidatorAddress(validator) {
		valAddress = validator
	} else {
		valAddress = l.GetAddress(validator)

		if !omg.IsValidatorAddress(valAddress) {
			return fmt.Errorf("invalid validator: %q", valAddress)
		}
	}

	// Check balance for delegator
	balanceBefore, err := omg.GetBalanceDec(delegatorAddress)
	if err != nil {
		return fmt.Errorf("error querying balance for %s", delegator)
	}
	r, err := omg.GetRewards(delegatorAddress)
	if err != nil {
		return err
	}
	if r.Total[0].Denom != cfg.Denom {
		return fmt.Errorf("expected total denom to be %q, got %q", cfg.Denom, r.Total[0].Denom)
	}
	rewards, err := omg.StrToDec(r.Total[0].Amount)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "Delegator         : %s [%s]\n", delegator, delegatorAddress)
	fmt.Fprintf(out, "Existing balance  : %s%s\n", omg.PrettifyDenom(balanceBefore), denom)
	fmt.Fprintf(out, "Unclaimed rewards : %s%s\n", omg.PrettifyDenom(rewards), denom)
	fmt.Fprintln(out, "----")
	fmt.Fprintf(out, "Withdrawing rewards...\n")
	omg.TxWithdrawRewards(out, delegator, keyring, auto)

	// Wait till balance is updated
	var balance decimal.Decimal
	count := 0
	for count <= 10 {
		balance, _ = omg.GetBalanceDec(delegatorAddress)
		if balance.GreaterThan(balanceBefore) {
			fmt.Fprintf(out, "...updated balance.\n")
			break
		}
		count++
		time.Sleep(1 * time.Second)
	}
	// If balance not updated and -auto flag set then abort
	if auto && balance == balanceBefore {
		return fmt.Errorf("balance not increased. Aborting auto-restake")
	}
	// Restake amount leaving approx remainder of 1 token
	amount := balance.Sub(remainAmt)
	fmt.Fprintln(out, "----")
	fmt.Fprintf(out, "Delegating to     : %s\n", valAddress)
	fmt.Fprintf(out, "Available balance : %s\n", omg.PrettifyDenom(balance))
	fmt.Fprintf(out, "Delegation amount : %s\n", omg.PrettifyDenom(amount))
	fmt.Fprintf(out, "Remander amount   : %s\n", omg.PrettifyDenom(remainAmt))
	fmt.Fprintln(out, "----")
	omg.TxDelegateToValidator(delegator, valAddress, amount, keyring, auto)
	return nil
}
