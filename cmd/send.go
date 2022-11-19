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

	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send [name] [name | address] [amount][denom]",
	Short: "Send tokens from an account to another account/address",
	Long:  `Send tokens from an account to another account/address`,
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
		err         error
	)
	l := &omg.Accounts{}

	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	if len(args) == 0 {
		from, to, err = omg.GetTxAccounts(os.Stdin, "send", args...)
		if err != nil {
			return err
		}
	} else if len(args) >= 2 {
		from = args[0]
		if omg.IsNormalAddress(args[1]) {
			to = args[1]
			toAddress = args[1]
		} else {
			to = args[1]
		}
	}
	fromAddress = l.GetAddress(from)
	if fromAddress == "" {
		return fmt.Errorf("No account found")
	}
	if !omg.IsNormalAddress(fromAddress) {
		return fmt.Errorf("Invalid from account: %s\n", fromAddress)
	}
	// Check if valid validator address
	if !omg.IsNormalAddress(to) {
		toAddress = l.GetAddress(to)
		if toAddress == "" {
			return fmt.Errorf("Address not in list\n")
		}
	}
	if !omg.IsNormalAddress(toAddress) {
		return fmt.Errorf("Invalid address: %s\n", toAddress)
	}
	amount, err := omg.GetAmount(os.Stdin, "send", fromAddress, args...)
	if err != nil {
		return err
	}
	// Check balance for sender
	if to == toAddress {
		fmt.Fprintf(out, "To    : %s\n", toAddress)
	} else {
		fmt.Fprintf(out, "To    : %s [%s]\n", to, toAddress)
	}
	fmt.Fprintf(out, "From  : %s [%s]\n", from, fromAddress)
	omg.CheckBalances(fromAddress)

	fmt.Fprintf(out, "Amount requested  : %s\n", omg.PrettifyAmount(amount, cfg.Denom))

	err = omg.TxSend(fromAddress, toAddress, amount, keyring, auto)
	if err != nil {
		return err
	}
	return nil
}
