/*
Copyright © 2022 dotneko
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"

	"github.com/spf13/cobra"
)

// delegateCmd represents the delegate command
var delegateCmd = &cobra.Command{
	Aliases: []string{"del", "d"},
	Use:     "delegate [account] [moniker|valoper-address] [amount][denom]",
	Short:   "Delegate tokens from account to validator",
	Long: `Delegate tokens from account to validator.
	
If no amount is provided, it is assumed that the user wants to delegate the full amount.
A remainder specified by the '--remainder' or ='-r' flag specifies the minimum estimated
remaining balance that must be left after delegation. The transaction will abort if there is
insufficient funds.
Therefore:

	[amount] must be >= [balance after withdraw rewards] - [remainder]

The remainder can be set in the configuration file.

Examples:

Delegate specified amount from user1 to validator1:
# omg tx delegate user1 validator1 1000000000anom

Delegate full balance (less default remainder):
# omg tx delegate user1 validator1

Delegate full balance and specify remainder:
# omg tx delegate user1 validator1 -r 1000000000anom
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		if err := cobra.RangeArgs(2, 3)(cmd, args); err != nil {
			return fmt.Errorf("expecting [account] [moniker|valoper-address] [amount][denom] as arguments")
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
		hash, err := cmd.Flags().GetBool("txhash")
		if err != nil {
			return err
		}
		var outType string = ""
		if hash {
			outType = omg.HASH
		}
		remainder, err := cmd.Flags().GetString("remainder")
		if err != nil {
			return err
		}
		return delegateAction(os.Stdout, auto, keyring, outType, remainder, args)
	},
}

func init() {
	txCmd.AddCommand(delegateCmd)

	delegateCmd.Flags().BoolP("full", "f", false, "Delegate full balance amount (less remainder)")
	delegateCmd.Flags().StringP("remainder", "r", cfg.Remainder, "Remainder after delegate")

}

func delegateAction(out io.Writer, auto bool, keyring, outType, remainder string, args []string) error {

	delegator := args[0]
	validator := args[1]
	var (
		delegatorAddress    string
		valAddress          string
		valoperMoniker      string
		normalizedAmt       string
		normalizedRemainder string
		amtCoin             sdktypes.Coin
		expectedBalance     sdktypes.Coin
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
			}
		}
	}
	// Check balance
	balance, err := omg.GetBalance(delegatorAddress)
	if err != nil {
		return fmt.Errorf("cannot get balance for %s: %s", delegatorAddress, err)
	}
	// Parse remainder
	normalizedRemainder, err = omg.NormalizeAmountDenom(remainder)
	if err != nil {
		return fmt.Errorf("cannot normalize remainder %s: %s", remainder, err)
	}
	remainderCoin, err := sdktypes.ParseCoinNormalized(normalizedRemainder)
	if err != nil {
		return fmt.Errorf("cannot parse remainder %s: %s", normalizedRemainder, err)
	}
	if len(args) < 3 {
		// Delegate full amount if none given
		amtCoin = balance.Sub(remainderCoin)
		expectedBalance = remainderCoin
	} else {
		// Parse specified delegation amount
		normalizedAmt, err = omg.NormalizeAmountDenom(args[2])
		if err != nil {
			return fmt.Errorf("cannot normalize amount %s: %s", args[2], err)
		}
		amtCoin, err = sdktypes.ParseCoinNormalized(normalizedAmt)
		if err != nil {
			return fmt.Errorf("cannot parse amount %s: %s", normalizedAmt, err)
		}
		expectedBalance = balance.Sub(amtCoin)
	}
	if amtCoin.IsNegative() || amtCoin.IsZero() {
		return fmt.Errorf("amount must be greater than zero, got %s", omg.PrettifyBaseAmt(amtCoin.String()))
	}
	if !balance.IsGTE(amtCoin.Add(remainderCoin)) {
		return fmt.Errorf("insufficient balance after deducting remainder: %s", omg.PrettifyBaseAmt(expectedBalance.String()))
	}
	if !(auto && outType == omg.HASH) {
		fmt.Fprintf(out, "Delegator             : %s [%s]\n", delegator, delegatorAddress)
		fmt.Fprintf(out, "Available balance     : %10s %s ( %s%s ) \n", omg.AmtToTokenStr(balance.String()), cfg.Token, omg.PrettifyBaseAmt(balance.String()), cfg.BaseDenom)
		fmt.Fprintf(out, "----\n")
		fmt.Fprintf(out, "Delegate to Validator : %s\n", valAddress)
		fmt.Fprintf(out, "Delegation amount     : %10s %s ( %s%s )\n", omg.AmtToTokenStr(amtCoin.String()), cfg.Token, omg.PrettifyBaseAmt(amtCoin.String()), cfg.BaseDenom)
		fmt.Fprintf(out, "Min remainder setting : %10s %s ( %s%s )\n", omg.AmtToTokenStr(remainderCoin.String()), cfg.Token, omg.PrettifyBaseAmt(remainderCoin.String()), cfg.BaseDenom)
		fmt.Fprintf(out, "Est minimum remaining : %10s %s ( %s%s )\n", omg.AmtToTokenStr(expectedBalance.String()), cfg.Token, omg.PrettifyBaseAmt(expectedBalance.String()), cfg.BaseDenom)
		fmt.Fprintf(out, "----\n")
	}
	txhash, err := omg.TxDelegateToValidator(out, delegator, valAddress, amtCoin.String(), auto, keyring, outType)
	if err != nil {
		return err
	}
	if outType == omg.HASH {
		fmt.Fprintln(out, txhash)
	}
	return nil
}
