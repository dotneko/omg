/*
Copyright © 2022 dotneko

*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"
	"github.com/shopspring/decimal"

	"github.com/spf13/cobra"
)

// delegateCmd represents the delegate command
var delegateCmd = &cobra.Command{
	Aliases: []string{"del", "d"},
	Use:     "delegate [account] [moniker|valoper-address] [amount][denom]",
	Short:   "Delegate tokens from account to validator",
	Long: `Delegate tokens from account to validator.
	
The '--full' or '-f' flag can be added to delegate the full available balance.
A remainder specified by the '--remainder' or ='-r' flag specifies the minimum estimated
remaining balance that must be left after delegation or the transaction will abort. 
Therefore:

	[amount] must be >= [balance after withdraw rewards] - [remainder]

The remainder can be set in the configuration file.

Examples:

Delegate specified amount from user1 to validator1:
# omg tx delegate user1 validator1 1000000000anom

Delegate full balance (less default remainder):
# omg tx delegate user1 validator1 --full

Delegate full balance and specify remainder:
# omg tx delegate user1 validator1 -f -r 1000000000anom
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		if err := cobra.RangeArgs(2, 3)(cmd, args); err != nil {
			return fmt.Errorf("expecting [account] [moniker|valoper-address] [amount][denom] as arguments")
		}
		all, _ := cmd.Flags().GetBool("full")
		if len(args) == 2 && !all {
			return fmt.Errorf("no delegation amount nor --full flag specified")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		keyring, err := cmd.Flags().GetString("keyring")
		if err != nil {
			return err
		}
		auto, err := cmd.Flags().GetBool("yes")
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
	delegateCmd.Flags().StringP("remainder", "r", cfg.Remainder, "Remainder after delegate")

}

func delegateAction(out io.Writer, keyring string, auto bool, all bool, remainder string, args []string) error {

	delegator := args[0]
	validator := args[1]
	var (
		delegatorAddress string
		valAddress       string
		valoperMoniker   string
		amount           decimal.Decimal
		denom            string = cfg.BaseDenom
		remainAmt        decimal.Decimal
		remainDenom      string
		expectedBalance  decimal.Decimal
	)

	l := &omg.Accounts{}
	if err := l.Load(cfg.OmgFilepath); err != nil {
		return err
	}
	// Check if delegator in list and is not validator account
	delegatorAddress = l.GetAddress(delegator)
	if delegatorAddress == "" {
		return fmt.Errorf("account %q not found", delegator)
	}
	if !omg.IsNormalAddress(delegatorAddress) {
		return fmt.Errorf("invalid delegator address for %s", delegator)
	}
	if delegatorAddress != "" && delegatorAddress != omg.QueryKeyringAddress(delegator, keyring) {
		return fmt.Errorf("delegator/address not in keyring")
	}
	// Check if valid validator or validator address or moniker
	if omg.IsValidatorAddress(validator) {
		valAddress = validator
	} else {
		valAddress = l.GetAddress(validator)
		if !omg.IsValidatorAddress(valAddress) {
			// Query chain for address matching moniker if not found in address book
			searchMoniker := strings.ToLower(validator)
			valoperMoniker, valAddress = omg.GetValidator(searchMoniker)
			if valoperMoniker == "" {
				return fmt.Errorf("no validator matching %s found", validator)
			} else {
				fmt.Fprintf(out, "Found active validator %s [%s]\n----\n", valoperMoniker, valAddress)
			}
		}
	}
	fmt.Fprintf(out, "Delegator             : %s [%s]\n", delegator, delegatorAddress)
	balance, err := omg.GetBalanceDec(delegatorAddress)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "Available balance     : %s%s\n", omg.PrettifyDenom(balance), cfg.BaseDenom)
	fmt.Fprintln(out, "----")
	fmt.Fprintf(out, "Delegate to Validator : %s\n", valAddress)

	// Parse remainder
	remainAmt, remainDenom, err = omg.StrSplitAmountDenomDec(remainder)
	if err != nil {
		return err
	}
	if remainDenom == cfg.Token {
		remainAmt, _ = omg.ConvertDecDenom(remainAmt, remainDenom)
	}
	if all {
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
	fmt.Fprintf(out, "Delegation amount     : %s%s\n", omg.PrettifyDenom(amount), cfg.BaseDenom)
	fmt.Fprintf(out, "Min remainder setting : %s%s\n", omg.PrettifyDenom(remainAmt), cfg.BaseDenom)
	if amount.IsNegative() || amount.IsZero() {
		return fmt.Errorf("amount must be greater than zero, got %s", omg.PrettifyDenom(amount))
	}
	if amount.GreaterThan(balance.Sub(remainAmt)) {
		return fmt.Errorf("insufficient balance after deducting remainder: %s %s", omg.PrettifyDenom(expectedBalance), denom)
	}
	fmt.Fprintf(out, "Est minimum remaining : %s%s\n", omg.PrettifyDenom(expectedBalance), cfg.BaseDenom)
	fmt.Fprintln(out, "----")
	omg.TxDelegateToValidator(delegator, valAddress, amount, keyring, auto)

	return nil
}
