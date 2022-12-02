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
var wdrewardsCmd = &cobra.Command{
	Aliases: []string{"wd", "wr", "w"},
	Use:     "withdraw-rewards [name]",
	Short:   "Withdraw all rewards for account",
	Long:    `Withdraw all rewards for account`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("expecting [name] as argument")
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
		return wdrewardsAction(os.Stdout, keyring, auto, args)
	},
}

func init() {
	txCmd.AddCommand(wdrewardsCmd)

}

func wdrewardsAction(out io.Writer, keyring string, auto bool, args []string) error {

	delegator := args[0]
	l := &omg.Accounts{}
	if err := l.Load(cfg.OmgFilepath); err != nil {
		return err
	}
	delegatorAddress := l.GetAddress(delegator)
	if delegatorAddress == "" {
		return fmt.Errorf("account %q not found", delegator)
	}
	if !omg.IsNormalAddress(delegatorAddress) {
		return fmt.Errorf("%s is not a valid account", delegator)
	}
	if delegatorAddress != "" && delegatorAddress != omg.QueryKeyringAddress(delegator, keyring) {
		return fmt.Errorf("delegator/address not in keyring")
	}
	err := omg.TxWithdrawRewards(out, delegator, keyring, auto)
	if err != nil {
		return err
	}
	return nil
}
