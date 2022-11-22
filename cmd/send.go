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
		auto, err := cmd.Flags().GetBool("auto")
		if err != nil {
			return err
		}
		keyring, err := cmd.Flags().GetString("keyring")
		if err != nil {
			return err
		}
		return sendAction(os.Stdout, auto, keyring, args)
	},
}

func init() {
	txCmd.AddCommand(sendCmd)

	sendCmd.Flags().BoolP("auto", "a", false, "Auto confirm transaction")
	sendCmd.Flags().StringP("keyring", "k", cfg.KeyringBackend, "Specify keyring-backend to use.")
}

func sendAction(out io.Writer, auto bool, keyring string, args []string) error {
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

	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	// Parse and validate [from]
	from = args[0]
	fromAddress = l.GetAddress(from)
	if fromAddress == "" {
		return fmt.Errorf("no account found")
	}
	if !omg.IsNormalAddress(fromAddress) {
		return fmt.Errorf("invalid from account: %s", fromAddress)
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
	if denom == cfg.Token {
		amount, denom = omg.ConvertDecDenom(amount, denom)
	}
	if err != nil {
		return err
	}
	// Check balance for sender
	balance, err := omg.GetBalanceDec(fromAddress)
	if err != nil {
		return fmt.Errorf("error querying balance for %s", from)
	}
	// Check balance for sender
	if to == toAddress {
		fmt.Fprintf(out, "To                : %s\n", toAddress)
	} else {
		fmt.Fprintf(out, "To                : %s [%s]\n", to, toAddress)
	}
	fmt.Fprintf(out, "From              : %s [%s]\n", from, fromAddress)
	fmt.Fprintf(out, "Available balance : %s\n", omg.PrettifyAmount(balance, denom))
	fmt.Fprintf(out, "Amount requested  : %s\n", omg.PrettifyAmount(amount, denom))
	fmt.Fprintln(out, "----")

	if amount.GreaterThan(balance) {
		return fmt.Errorf("insufficient balance for send amount")
	}
	err = omg.TxSend(fromAddress, toAddress, amount, keyring, auto)
	if err != nil {
		return err
	}
	return nil
}
