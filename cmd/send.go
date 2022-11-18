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
	Use:   "send [from-account alias] [to-alias] [amount][denom]",
	Short: "send [from-account alias] [to-alias] [amount][denom]",
	Long:  `Send tokens between two accounts in the address book`,
	RunE: func(cmd *cobra.Command, args []string) error {
		auto, err := cmd.Flags().GetBool("auto")
		if err != nil {
			return err
		}
		keyring, err := cmd.Flags().GetString("keyring")
		return sendAction(os.Stdout, auto, keyring, args)
	},
}

func init() {
	txCmd.AddCommand(sendCmd)

	sendCmd.Flags().BoolP("auto", "a", false, "Auto confirm transaction")
	sendCmd.Flags().StringP("keyring", "k", cfg.KeyringBackend, "Specify keyring-backend to use.")
}

func sendAction(out io.Writer, auto bool, keyring string, args []string) error {
	from, to, err := omg.GetTxAccounts(os.Stdin, "send", args...)
	if err != nil {
		return err
	}
	// Check if delegator in list and is not validator account
	l := &omg.Accounts{}

	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	fromAddress := l.GetAddress(from)
	if fromAddress == "" {
		return fmt.Errorf("Error: no from address")
	}
	if !omg.IsNormalAddress(fromAddress) {
		return fmt.Errorf("Error: Invalid normal account: %s\n", fromAddress)
	}
	// Check if valid validator address
	toAddress := l.GetAddress(to)
	if toAddress == "" {
		return fmt.Errorf("Error: Address not in list\n")
	}
	if !omg.IsNormalAddress(toAddress) {
		return fmt.Errorf("Error: Invalid normal account: %s\n", toAddress)
	}
	amount, err := omg.GetAmount(os.Stdin, "send", fromAddress, args...)
	if err != nil {
		return err
	}
	// Check balance for delegator
	fmt.Fprintf(out, "To   : %s [%s]\n", to, toAddress)
	fmt.Fprintf(out, "From : %s [%s]\n", from, fromAddress)
	omg.CheckBalances(fromAddress)

	fmt.Fprintf(out, "Amount requested: %s\n", omg.PrettifyAmount(amount, cfg.Denom))

	err = omg.TxSend(fromAddress, toAddress, amount, keyring, auto)
	if err != nil {
		return err
	}
	return nil
}
