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
	Aliases: []string{"r"},
	Use:     "restake [name] [validator|valoper-address] *OPTIONAL:[amount][denom]",
	Short:   "Withdraw rewards and restake to validator",
	Long: fmt.Sprintf(`Withdraw all rewards for account, then re-delegate to validator.

A remainder specified by the '--remainder' or ='-f' flag will be deducted from the 
updated balance after rewards withdrawn, and is effective even if a delegation amount
is specified.

If a delegation amount is specified, the final balance after delegation must exceed the
remainder or the transaction will abort. Therefore:

	[amount] must be >= [balance after withdraw rewards] - [remainder]

The remainder can be set in the configuration file, currently: %s

Examples:

Restake full balance (less default remainder):
# omg tx restake user1 validator1

Restake full balance (specify remainder):
# omg tx restake user1 validator1 -r 1000000anom

Restake specified amount
# omg tx restake user1 validator1 1nom
	`, cfg.Remainder),
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		if err := cobra.RangeArgs(2, 3)(cmd, args); err != nil {
			fmt.Printf("Error: %s\n", fmt.Errorf("expecting [account] [validator|valoper-address] as arguments"))
		}
		return nil
	},
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

	restakeCmd.Flags().StringP("remainder", "r", cfg.Remainder, "Remainder after restake")

}

func restakeAction(out io.Writer, remainder string, keyring string, auto bool, args []string) error {
	// Ensure all arguments provided

	delegator := args[0]
	validator := args[1]
	var (
		delegatorAddress string
		valAddress       string
		amount           decimal.Decimal
		balanceBefore    decimal.Decimal
		denom            string = cfg.BaseDenom
		remainAmt        decimal.Decimal
		remainDenom      string
		expectedBalance  decimal.Decimal
	)

	l := &omg.Accounts{}
	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	// Check if delegator in list and is not validator account
	delegatorAddress = l.GetAddress(delegator)

	if !omg.IsNormalAddress(delegatorAddress) {
		return fmt.Errorf("invalid delegator address for %s", delegator)
	}
	// Check if valid validator or validator address
	if omg.IsValidatorAddress(validator) {
		valAddress = validator
	} else {
		valAddress = l.GetAddress(validator)
		if !omg.IsValidatorAddress(valAddress) {
			return fmt.Errorf("invalid validator address %s", valAddress)
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
	if r.Total[0].Denom != cfg.BaseDenom {
		return fmt.Errorf("expected total denom to be %q, got %q", cfg.BaseDenom, r.Total[0].Denom)
	}
	fmt.Fprintf(out, r.Total[0].Amount)
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
	// Parse remainder
	remainAmt, remainDenom, err = omg.StrSplitAmountDenomDec(remainder)
	if err != nil {
		return err
	}
	if remainDenom == cfg.Token {
		remainAmt, _ = omg.ConvertDecDenom(remainAmt, remainDenom)
	}
	if len(args) < 3 {
		// Restake full balance less remainder if no amount specified
		amount = balance.Sub(remainAmt)
		expectedBalance = remainAmt
	} else {
		// Parse delegation amount
		amount, denom, err = omg.StrSplitAmountDenomDec(args[2])
		if err != nil {
			return err
		}
		// Convert to baseDenom if denominated in Token
		if denom == cfg.Token {
			amount, denom = omg.ConvertDecDenom(amount, denom)
		}
		expectedBalance = balance.Sub(amount)
	}
	fmt.Fprintln(out, "----")
	fmt.Fprintf(out, "Delegate to Validator : %s\n", valAddress)
	fmt.Fprintf(out, "Available balance     : %s%s\n", omg.PrettifyDenom(balance), cfg.BaseDenom)
	fmt.Fprintf(out, "Delegation amount     : %s%s\n", omg.PrettifyDenom(amount), cfg.BaseDenom)
	fmt.Fprintf(out, "Remainder amount      : %s%s\n", omg.PrettifyDenom(remainAmt), cfg.BaseDenom)
	if amount.GreaterThan(balance.Sub(remainAmt)) {
		return fmt.Errorf("insufficient balance after deducting remainder: %s %s", amount.String(), denom)
	}
	fmt.Fprintf(out, "Expected after Tx     : %s%s\n", omg.PrettifyDenom(expectedBalance), cfg.BaseDenom)
	fmt.Fprintln(out, "----")

	omg.TxDelegateToValidator(delegator, valAddress, amount, keyring, auto)
	return nil
}
