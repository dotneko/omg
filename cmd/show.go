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
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.RangeArgs(0, 1)(cmd, args); err != nil {
			cmd.Help()
			os.Exit(0)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		filterNormal, err := cmd.Flags().GetBool("regular")
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
	addressCmd.AddCommand(showCmd)

	showCmd.Flags().BoolP("address", "a", false, "Show addresses only")
	showCmd.Flags().BoolP("regular", "r", false, "Select regular accounts")
	showCmd.Flags().BoolP("validator", "v", false, "Select validator accounts")
}

func showAction(out io.Writer, filterAccount string, addressOnly bool, args []string) error {
	l := &omg.Accounts{}

	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	if len(*l) == 0 {
		return fmt.Errorf("no accounts in store")
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
			return fmt.Errorf("no address found for %s", args[0])
		}
		fmt.Fprint(out, address)
	}
	return nil
}
