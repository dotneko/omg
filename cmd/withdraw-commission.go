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
		keyring, err := cmd.Flags().GetString("keyring")
		if err != nil {
			return err
		}
		auto, err := cmd.Flags().GetBool("auto")
		if err != nil {
			return err
		}
		return wdCommissionAction(os.Stdout, keyring, auto, args)
	},
}

func init() {
	txCmd.AddCommand(wdCommissionCmd)

}

func wdCommissionAction(out io.Writer, keyring string, auto bool, args []string) error {
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
	err := omg.TxWithdrawValidatorCommission(out, delegator, valoperAddress, keyring, auto)
	if err != nil {
		return err
	}
	return nil
}
