/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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

// rewardsCmd represents the rewards command
var rewardsCmd = &cobra.Command{
	Aliases: []string{"rw", "r"},
	Use:     "rewards [name|address]",
	Short:   "Query rewards for an account or address",
	Long:    `Query rewards for an account or address.`,
	Args: func(cmd *cobra.Command, args []string) error {

		if err := cobra.RangeArgs(0, 1)(cmd, args); err != nil {
			return fmt.Errorf("expecting [name|address] as argument or --all flag")
		}
		allAccounts, _ := cmd.Flags().GetBool("all")
		if len(args) == 0 && !allAccounts {
			cmd.Help()
			os.Exit(0)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		allAccounts, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}
		detail, err := cmd.Flags().GetBool("detail")
		if err != nil {
			return err
		}
		raw, err := cmd.Flags().GetBool("raw")
		if err != nil {
			return err
		}
		return rewardsAction(os.Stdout, allAccounts, detail, raw, args)
	},
}

func init() {
	rootCmd.AddCommand(rewardsCmd)

	rewardsCmd.Flags().BoolP("all", "a", false, "Check all accounts in address book")
	rewardsCmd.Flags().BoolP("detail", "d", false, "Detailed output")
	rewardsCmd.Flags().BoolP("raw", "r", false, "Raw output")

}

func rewardsAction(out io.Writer, allAccounts bool, detail bool, raw bool, args []string) error {
	l := &omg.Accounts{}

	if err := l.Load(cfg.OmgFilepath); err != nil {
		return err
	}

	var outType string = ""
	if raw {
		outType = omg.RAW
	} else if detail {
		outType = omg.DETAIL
	}
	if allAccounts {
		for _, acc := range *l {
			if omg.IsNormalAddress(acc.Address) {
				r, err := omg.GetRewards(acc.Address)
				if err != nil {
					return err
				}
				if len(r.Rewards) != 0 && len(r.Rewards[0].Reward) != 0 {
					fmt.Printf("%s:\n", acc.Alias)
					for _, v := range r.Rewards {
						amt, err := omg.StrToDec(v.Reward[0].Amount)
						if err != nil {
							return err
						}
						moniker, valoperAddress := omg.GetValidator(v.ValidatorAddress)
						omg.OutputAmount(out, moniker, valoperAddress, amt, v.Reward[0].Denom, outType)
					}
				}
			}
		}
		return nil
	}
	if len(args) == 0 {
		return fmt.Errorf("no account/address provided")
	}
	var address string
	if omg.IsNormalAddress(args[0]) {
		address = args[0]
	} else {
		address = l.GetAddress(args[0])
		if address == "" {
			return fmt.Errorf("account %q not found", args[0])
		}
	}
	r, err := omg.GetRewards(address)
	if err != nil {
		return err
	}
	for _, v := range r.Rewards {
		amt, err := omg.StrToDec(v.Reward[0].Amount)
		if err != nil {
			return err
		}
		moniker, valoperAddress := omg.GetValidator(v.ValidatorAddress)
		omg.OutputAmount(out, moniker, valoperAddress, amt, v.Reward[0].Denom, outType)
	}
	return nil
}
