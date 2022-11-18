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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Aliases: []string{"l"},
	Use:     "list",
	Short:   "List addresses",
	Long:    `Lists addresses and their aliases`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filterNormal, err := cmd.Flags().GetBool("normal")
		if err != nil {
			return err
		}
		filterValoper, err := cmd.Flags().GetBool("validator")
		if err != nil {
			return err
		}
		addressOnly, err := cmd.Flags().GetBool("address")
		if err != nil {
			return err
		}
		filterAddrType := ""
		if filterNormal && !filterValoper {
			filterAddrType = omg.AccNormal
		} else if filterValoper && !filterNormal {
			filterAddrType = omg.AccValoper
		}
		return listAction(os.Stdout, filterAddrType, addressOnly, args)
	},
}

func init() {
	addrCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("address", "a", false, "List addresses only")
	listCmd.Flags().BoolP("normal", "n", false, "Select normal accounts")
	listCmd.Flags().BoolP("validator", "v", false, "Select validator accounts")
}

func listAction(out io.Writer, filterAccount string, addressOnly bool, args []string) error {
	l := &omg.Accounts{}

	// Read from saved address book
	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	if len(*l) == 0 {
		return fmt.Errorf("No accounts in store")
	}
	if filterAccount == "" && !addressOnly {
		fmt.Fprint(out, l)
		return nil
	}
	fmt.Fprint(out, l.ListFiltered(filterAccount, addressOnly))
	return nil
}
