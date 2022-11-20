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
	Aliases: []string{"wdall", "wd", "w"},
	Use:     "wdrewards [alias]",
	Short:   "Withdraw all rewards",
	Long:    `Withdraw all rewards for account`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyring, err := cmd.Flags().GetString("keyring")
		if err != nil {
			return err
		}
		auto, err := cmd.Flags().GetBool("auto")
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

	l := &omg.Accounts{}
	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	address := l.GetAddress(args[0])
	if address == "" {
		return fmt.Errorf("Account %q not found.\n", args[0])
	}
	if !omg.IsNormalAddress(address) {
		return fmt.Errorf("%s is not a normal account", args[0])
	}
	err := omg.TxWithdrawRewards(out, args[0], keyring, auto)
	if err != nil {
		return err
	}
	return nil
}
