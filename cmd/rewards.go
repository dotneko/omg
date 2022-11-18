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
	Use:     "rewards [alias]",
	Short:   "rewards [alias] or 'rewards -a' for all accounts",
	Long:    `Query rewards for an account.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		allAccounts, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}
		return rewardsAction(os.Stdout, allAccounts, args)
	},
}

func init() {
	queryCmd.AddCommand(rewardsCmd)

	rewardsCmd.Flags().BoolP("all", "a", false, "Check all accounts in address book")
}

func rewardsAction(out io.Writer, allAccounts bool, args []string) error {
	l := &omg.Accounts{}

	if err := l.Load(cfg.OmgFilename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if allAccounts {
		for _, acc := range *l {
			if omg.IsNormalAddress(acc.Address) {
				r, err := omg.GetRewards(acc.Address)
				if err != nil {
					return err
				}
				fmt.Printf("Rewards for %10s [%s]:\n", acc.Alias, acc.Address)
				for _, v := range r.Rewards {
					amt, err := omg.StrToFloat(v.Reward[0].Amount)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Printf(" - %s - %8.5f %s\n", v.ValidatorAddress, omg.DenomToToken(amt), cfg.Token)
				}
				fmt.Println()
			}
		}
		return nil
	}
	if len(args) == 0 {
		return fmt.Errorf("No account provided.\n")
	}
	address := l.GetAddress(args[0])
	if address == "" {
		return fmt.Errorf("Aaccount %q not found.\n", args[0])
	}
	r, err := omg.GetRewards(address)
	if err != nil {
		return err
	}
	for _, v := range r.Rewards {
		amt, err := omg.StrToFloat(v.Reward[0].Amount)
		if err != nil {
			return err
		}
		fmt.Printf("%s - %8.5f nom\n", v.ValidatorAddress, omg.DenomToToken(amt))
	}
	return nil
}
