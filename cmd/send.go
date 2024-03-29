/*
Copyright © 2022 dotneko
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

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Aliases: []string{"s"},
	Use:     "send [from: name] [to: name|address] [amount][denom]",
	Short:   "Send tokens from an account to another account/address",
	Long:    `Send tokens from an account to another account/address`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		if err := cobra.ExactArgs(3)(cmd, args); err != nil {
			return fmt.Errorf("expecting [from: name] [to: name|address] [amount][denom] as arguments")
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
		return sendAction(os.Stdout, auto, keyring, outType, args)
	},
}

func init() {
	txCmd.AddCommand(sendCmd)

}

func sendAction(out io.Writer, auto bool, keyring, outType string, args []string) error {
	var (
		from        string
		to          string
		fromAddress string
		toAddress   string
		amtCoin     sdktypes.Coin
		err         error
	)
	l := &omg.Accounts{}

	if err := l.Load(cfg.OmgFilepath); err != nil {
		return err
	}
	// Parse and validate [from]
	from = args[0]
	fromAddress = l.GetAddress(from)
	if fromAddress == "" {
		return fmt.Errorf("account %q not found", from)
	}
	if !omg.IsNormalAddress(fromAddress) {
		return fmt.Errorf("invalid from account: %s", fromAddress)
	}
	if fromAddress != "" && fromAddress != omg.QueryKeyringAddress(from, keyring) {
		return fmt.Errorf("delegator/address not in keyring")
	}

	// Parse and validate [to]
	if omg.IsNormalAddress(args[1]) {
		to = args[1]
		toAddress = args[1]
	} else {
		to = args[1]
	}
	if !omg.IsNormalAddress(to) {
		toAddress = l.GetAddress(to)
		if toAddress == "" {
			return fmt.Errorf("address not in list")
		}
	}
	if !omg.IsNormalAddress(toAddress) {
		return fmt.Errorf("invalid address: %s", toAddress)
	}
	// Normalize denom capitalization
	normalizedAmt, err := omg.NormalizeAmountDenom(args[2])
	if err != nil {
		return fmt.Errorf("failed to normalize %s: %s", args[2], err)
	}
	amtCoin, err = sdktypes.ParseCoinNormalized(normalizedAmt)
	if err != nil {
		return fmt.Errorf("failed to parse %s to Coin: %s", normalizedAmt, err)
	}
	// Check balance for sender
	balance, err := omg.GetBalance(fromAddress)
	if err != nil {
		return fmt.Errorf("error querying balance for %s", from)
	}
	balanceToken, _ := omg.AmtToTokenDecCoin(balance.String())
	requestAmtToken, err := omg.AmtToTokenDecCoin(amtCoin.String())
	if err != nil {
		return fmt.Errorf("error converting request amount: %s", err)
	}
	// Display transaction summary
	if !(auto && outType == omg.HASH) {
		if to == toAddress {
			fmt.Fprintf(out, "To                : %s\n", toAddress)
		} else {
			fmt.Fprintf(out, "To                : %s [%s]\n", to, toAddress)
		}
		fmt.Fprintf(out, "From              : %s [%s]\n", from, fromAddress)

		fmt.Fprintf(out, "Available balance : %s %s ( %s%s )\n", balanceToken, cfg.Token, omg.PrettifyBaseAmt(balance.String()), cfg.BaseDenom)
		fmt.Fprintf(out, "Amount requested  : %s %s ( %s )\n", requestAmtToken, cfg.Token, omg.PrettifyBaseAmt(amtCoin.String()))
		fmt.Fprintf(out, "----\n")
	}

	if amtCoin.IsGTE(balance) {
		return fmt.Errorf("insufficient balance for send amount")
	}
	txhash, err := omg.TxSend(out, fromAddress, toAddress, amtCoin.String(), auto, keyring, outType)
	if err != nil {
		return err
	}
	if outType == omg.HASH && txhash != "" {
		fmt.Fprintln(out, txhash)
	}
	return nil
}
