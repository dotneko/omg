/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"

	"github.com/spf13/cobra"
)

// rewardsCmd represents the rewards command
var rewardsCmd = &cobra.Command{
	Use:   "rewards [alias]",
	Short: "Query rewards for an account",
	Long:  `Query rewards for an account.`,
	Run: func(cmd *cobra.Command, args []string) {
		l := &omg.Accounts{}

		if err := l.Load(cfg.OmgFilename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		address := l.GetAddress(args[0])
		if address == "" {
			fmt.Printf("Error: account %q not found.\n", args[0])
			os.Exit(1)
		}
		r, err := omg.GetRewards(address)
		if err != nil {
			fmt.Sprintln(err)
		}
		for _, v := range r.Rewards {
			amt, err := omg.StrToFloat(v.Reward[0].Amount)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%s - %8.5f nom\n", v.ValidatorAddress, omg.DenomToToken(amt))
		}
	},
}

func init() {
	queryCmd.AddCommand(rewardsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rewardsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rewardsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
