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

// showCmd represents the show command
var showCmd = &cobra.Command{
	Aliases: []string{"list", "get", "s"},
	Use:     "show [name]",
	Short:   "Show one or all addresses",
	Long:    `Show one or more addresses for a given name in the address book`,
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
		return showAction(os.Stdout, filterAddrType, addressOnly, args)
	},
}

func init() {
	addrCmd.AddCommand(showCmd)

	showCmd.Flags().BoolP("address", "a", false, "show addresses only")
	showCmd.Flags().BoolP("normal", "n", false, "Select normal accounts")
	showCmd.Flags().BoolP("validator", "v", false, "Select validator accounts")
}

func showAction(out io.Writer, filterAccount string, addressOnly bool, args []string) error {
	l := &omg.Accounts{}

	// Read from saved address book
	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	if len(*l) == 0 {
		return fmt.Errorf("No accounts in store")
	}
	// List all addresses as default
	if len(args) == 0 {
		if filterAccount == "" && !addressOnly {
			fmt.Fprint(out, l)
			return nil
		}
		fmt.Fprint(out, l.ListFiltered(filterAccount, addressOnly))
		return nil
	}
	if len(args) == 1 {
		address := l.GetAddress(args[0])
		if address == "" {
			return fmt.Errorf("No address found for %s", args[0])
		}
		fmt.Fprint(out, address)
	}
	return nil
}
