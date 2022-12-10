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

// wdrewardsCmd represents the wdrewards command
var wdCommissionCmd = &cobra.Command{
	Aliases: []string{"commission", "wvc", "wc"},
	Use:     "withdraw-commission [name] [moniker|valoper-address]",
	Short:   "Withdraw commissions and rewards for validator",
	Long:    `Withdraw commissions and rewards for validator. Assumes the `,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		if err := cobra.ExactArgs(2)(cmd, args); err != nil {
			return fmt.Errorf("expecting [name] [moniker|valoper-address] as argument")
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
		return wdCommissionAction(os.Stdout, auto, keyring, outType, args)
	},
}

func init() {
	txCmd.AddCommand(wdCommissionCmd)

}

func wdCommissionAction(out io.Writer, auto bool, keyring, outType string, args []string) error {
	var (
		delegator        string = ""
		delegatorAddress string = ""
		valoperAddress   string = ""
	)
	delegator = args[0]

	// Check if valid account
	l := &omg.Accounts{}
	if err := l.Load(cfg.OmgFilepath); err != nil {
		return err
	}
	delegatorAddress = l.GetAddress(delegator)
	if delegatorAddress == "" {
		return fmt.Errorf("account %q not found", delegator)
	}
	if !omg.IsNormalAddress(delegatorAddress) {
		return fmt.Errorf("%s is not a valid account", delegator)
	}
	if delegatorAddress != "" && delegatorAddress != omg.QueryKeyringAddress(delegator, keyring) {
		return fmt.Errorf("delegator/address not in keyring")
	}

	// Check if valid moniker/valoper-address
	_, valoperAddress = omg.GetValidator(args[1])
	if valoperAddress == "" {
		return fmt.Errorf("%s is not a valid moniker/valoper-address", args[1])
	}
	txhash, err := omg.TxWithdrawValidatorCommission(out, delegator, valoperAddress, auto, keyring, outType)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	if outType == omg.HASH && txhash != "" {
		fmt.Fprintln(out, txhash)
	}
	return nil
}
