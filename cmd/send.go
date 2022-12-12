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
		amount      decimal.Decimal
		denom       string
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
	// Parse amount
	amount, denom, err = omg.StrSplitAmountDenomDec(args[2])
	if strings.EqualFold(denom, cfg.Token) {
		amount, _ = omg.ConvertDecDenom(amount, denom)
	} else if !strings.EqualFold(denom, cfg.BaseDenom) {
		return fmt.Errorf("unexpected denom, aborting")
	}
	if err != nil {
		return err
	}
	// Check balance for sender
	balance, err := omg.GetBalanceDec(fromAddress)
	if err != nil {
		return fmt.Errorf("error querying balance for %s", from)
	}
	// Display transaction summary
	if outType != omg.HASH {
		if to == toAddress {
			fmt.Fprintf(out, "To                : %s\n", toAddress)
		} else {
			fmt.Fprintf(out, "To                : %s [%s]\n", to, toAddress)
		}
		fmt.Fprintf(out, "From              : %s [%s]\n", from, fromAddress)
		fmt.Fprintf(out, "Available balance : %s %s ( %s%s )\n", omg.DenomToTokenDec(balance).StringFixed(4), cfg.Token, omg.PrettifyDenom(balance), cfg.BaseDenom)
		fmt.Fprintf(out, "Amount requested  : %s %s ( %s%s )\n", omg.DenomToTokenDec(amount).StringFixed(4), cfg.Token, omg.PrettifyDenom(amount), cfg.BaseDenom)
		fmt.Fprintf(out, "----\n")
	}

	if amount.GreaterThan(balance) {
		return fmt.Errorf("insufficient balance for send amount")
	}
	txhash, err := omg.TxSend(out, fromAddress, toAddress, amount, auto, keyring, outType)
	if err != nil {
		return err
	}
	if outType == omg.HASH && txhash != "" {
		fmt.Fprintln(out, txhash)
	}
	return nil
}
