/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"
	"github.com/spf13/cobra"
)

// restakeCmd represents the restake command
var restakeCmd = &cobra.Command{
	Aliases: []string{"r"},
	Use:     "restake [name] [moniker|valoper-address] *OPTIONAL:[amount][denom]",
	Short:   "Withdraw rewards and restake to validator",
	Long: fmt.Sprintf(`Withdraw all rewards for account, then re-delegate to validator.

A remainder specified by the '--remainder' or ='-r' flag specifies the minimum estimated
remaining balance that must be left after delegation or the transaction will abort. Therefore:

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
			return fmt.Errorf("expecting [account] [moniker|valoper-address] as arguments")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		auto, err := cmd.Flags().GetBool("yes")
		if err != nil {
			return err
		}
		keyring, err := cmd.Flags().GetString("keyring")
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
		withCommission, err := cmd.Flags().GetBool("commission")
		if err != nil {
			return err
		}
		remainder, err := cmd.Flags().GetString("remainder")
		if err != nil {
			return err
		}
		return restakeAction(os.Stdout, auto, keyring, outType, remainder, withCommission, args)
	},
}

func init() {
	txCmd.AddCommand(restakeCmd)
	restakeCmd.Flags().BoolP("commission", "c", false, "Include commission if validator")
	restakeCmd.Flags().StringP("remainder", "r", cfg.Remainder, "Remainder after restake")

}

func restakeAction(out io.Writer, auto bool, keyring, outType, remainder string, withCommission bool, args []string) error {
	// Ensure all arguments provided

	delegator := args[0]
	validator := args[1]
	var (
		delegatorAddress string
		valoperAddress   string
		valoperMoniker   string
		amount           sdktypes.Coin
		balanceBefore    sdktypes.Coin
		denom            string = cfg.BaseDenom
		remainderCoin    sdktypes.Coin
		rewards          sdktypes.Coin
		commission       sdktypes.Coin
		expectedBalance  sdktypes.Coin
		wdtxhash         string
		delegtxhash      string
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
	if delegatorAddress != omg.QueryKeyringAddress(delegator, keyring) {
		return fmt.Errorf("delegator/address not in keyring")
	}

	// Check if valid validator or validator address or moniker
	if omg.IsValidatorAddress(validator) {
		valoperAddress = validator
	} else {
		valoperAddress = l.GetAddress(validator)
		if !omg.IsValidatorAddress(valoperAddress) {
			// Query chain for address matching moniker if not found in address book
			searchMoniker := strings.ToLower(validator)
			valoperMoniker, valoperAddress = omg.GetValidator(searchMoniker)
			if valoperMoniker == "" {
				return fmt.Errorf("no validator matching %s found", validator)
			} else {
				fmt.Fprintf(out, "Found active validator %s [%s]\n----\n", valoperMoniker, valoperAddress)
			}
		}
	}
	// If include commissions, check if self-delegate
	if withCommission {
		if !omg.IsSelfDelegate(delegatorAddress, valoperAddress) {
			return fmt.Errorf("cannot include commissions: %s is not a self-delegate for %s", delegator, validator)
		}
		c, err := omg.QueryCommission(valoperAddress)
		if err != nil {
			return err
		}
		if c.Commission[0].Denom != cfg.BaseDenom {
			return fmt.Errorf("expected total denom to be %q, got %q", cfg.BaseDenom, c.Commission[0].Denom)
		}
		commission, err = sdktypes.ParseCoinNormalized(c.Commission[0].Amount + c.Commission[0].Denom)
		if err != nil {
			return err
		}
	}

	// Check balance for delegator
	balanceBefore, err := omg.GetBalance(delegatorAddress)
	if err != nil {
		return fmt.Errorf("error querying balance for %s", delegator)
	}

	// Check rewards
	r, err := omg.GetRewards(delegatorAddress)
	if err != nil {
		return err
	}
	if r.Total[0].Denom != cfg.BaseDenom {
		return fmt.Errorf("expected total denom to be %q, got %q", cfg.BaseDenom, r.Total[0].Denom)
	}
	rewards, err = sdktypes.ParseCoinNormalized(r.Total[0].Amount + r.Total[0].Denom)
	if err != nil {
		return err
	}
	if !(auto && outType == omg.HASH) {
		fmt.Fprintf(out, "Delegator             : %s [%s]\n", delegator, delegatorAddress)
		fmt.Fprintf(out, "Existing balance      : %s ( %s )\n", omg.AmtToTokenStr(balanceBefore.String()), omg.PrettifyBaseAmt(balanceBefore.String()))
		fmt.Fprintf(out, "Unclaimed rewards     : %s ( %s )\n", omg.AmtToTokenStr(rewards.String()), omg.PrettifyBaseAmt(rewards.String()))
		if withCommission {
			fmt.Fprintf(out, "Unclaimed commissions : %10s ( %s )\n", omg.AmtToTokenStr(commission.String()), omg.PrettifyBaseAmt(commission.String()))
			fmt.Fprintf(out, "----\n")
		}
	}
	if withCommission {
		if outType != omg.HASH {
			fmt.Fprintf(out, "Withdrawing rewards plus commission...\n")
		}
		wdtxhash, err = omg.TxWithdrawValidatorCommission(out, delegator, valoperAddress, auto, keyring, outType)
		if err != nil {
			return err
		}
	} else {
		if outType != omg.HASH {
			fmt.Fprintf(out, "Withdrawing rewards...\n")
		}
		wdtxhash, err = omg.TxWithdrawRewards(out, delegator, auto, keyring, outType)
		if err != nil {
			return err
		}
	}
	if outType == omg.HASH {
		fmt.Fprintf(out, "Withdraw hash: %s", wdtxhash)
	}
	// Wait till balance is updated
	var balance sdktypes.Coin
	count := 0
	if outType != omg.HASH {
		fmt.Fprintf(out, "Checking balance")
	}
	for count <= 10 {
		if outType != omg.HASH {
			fmt.Fprintf(out, ".")
		}
		balance, _ = omg.GetBalance(delegatorAddress)
		if !balance.IsLT(balanceBefore) {
			if outType != omg.HASH {
				fmt.Fprintf(out, "...updated.\n")
			}
			break
		}
		count++
		time.Sleep(3 * time.Second)
	}
	if outType != omg.HASH {
		fmt.Fprintf(out, "\n")
	}
	// if timeout
	if count > 10 {
		return fmt.Errorf("timed out retrieving balance. Aborting auto-restake")
	}
	// If balance not updated and -auto flag set then abort
	if auto && balance == balanceBefore {
		return fmt.Errorf("balance not increased. Aborting auto-restake")
	}
	// Parse remainder
	remainderCoin, err = sdktypes.ParseCoinNormalized(remainder)
	if err != nil {
		return err
	}
	if len(args) < 3 {
		// Restake full balance less remainder if no amount specified
		amount = balance.Sub(remainderCoin)
		expectedBalance = remainderCoin
	} else {
		// Parse delegation amount
		amount, err = sdktypes.ParseCoinNormalized(args[2])
		if err != nil {
			return err
		}
		expectedBalance = balance.Sub(amount)
	}
	if amount.IsNegative() || amount.IsZero() {
		return fmt.Errorf("amount must be greater than zero, got %s", omg.PrettifyBaseAmt(amount.String()))
	}
	if !amount.IsLT(balance.Sub(remainderCoin)) {
		return fmt.Errorf("insufficient balance after deducting remainder: %s %s", omg.PrettifyBaseAmt(expectedBalance.String()), denom)
	}
	if outType != omg.HASH {
		fmt.Fprintf(out, "----\n")
		fmt.Fprintf(out, "Delegate to Validator : %s\n", valoperAddress)
		fmt.Fprintf(out, "Available balance     : %s ( %s )\n", omg.AmtToTokenStr(balance.String()), omg.PrettifyBaseAmt(balance.String()))
		fmt.Fprintf(out, "Delegation amount     : %s ( %s )\n", omg.AmtToTokenStr(amount.String()), omg.PrettifyBaseAmt(amount.String()))
		fmt.Fprintf(out, "Min remainder setting : %s ( %s )\n", omg.AmtToTokenStr(remainderCoin.String()), omg.PrettifyBaseAmt(remainderCoin.String()))
		fmt.Fprintf(out, "Est minimum remaining : %s ( %s )\n", omg.AmtToTokenStr(expectedBalance.String()), omg.PrettifyBaseAmt(expectedBalance.String()))
		fmt.Fprint(out, "----\n")
	}

	delegtxhash, err = omg.TxDelegateToValidator(out, delegator, valoperAddress, amount.String(), auto, keyring, outType)
	if err != nil {
		return err
	}
	if outType == omg.HASH {
		fmt.Fprintf(out, "Delegate hash: %s", delegtxhash)
	}
	return nil
}
