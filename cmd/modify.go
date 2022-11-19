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

// modifyCmd represents the modify command
var modifyCmd = &cobra.Command{
	Aliases: []string{"mod", "mv", "rename"},
	Use:     "modify [name] [new name]",
	Short:   "modify [name] [new name] *optional: -a [address]",
	Long: `Modify an address book entry. Examples:

modify user1 user2
modify user1 user2 -a newaddress1234567890123456789
modify user1 -a newaddress1234567890123456789`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("address")
		if err != nil {
			return err
		}
		if len(args) == 1 && address == "" {
			return fmt.Errorf("Missing new name or address flag")
		}
		return modifyAction(os.Stdout, address, args)
	},
}

func init() {
	addrCmd.AddCommand(modifyCmd)

	modifyCmd.Flags().StringP("address", "a", "", "New address to use")
}

func modifyAction(out io.Writer, address string, args []string) error {

	var newAlias string
	oldAlias := args[0]
	if len(args) > 1 {
		newAlias = args[1]
	} else if address != "" {
		newAlias = oldAlias
	}

	l := &omg.Accounts{}

	// Read from saved address book
	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	if len(*l) == 0 {
		return fmt.Errorf("No accounts in store")
	}

	if len(newAlias) < cfg.MinAliasLength {
		return fmt.Errorf("Please use alias of at least 3 characters")
	}
	idx := l.GetIndex(oldAlias)
	err := l.Modify(idx, newAlias, address)
	if err != nil {
		return err
	}
	err = l.Save(cfg.OmgFilename)
	if err != nil {
		return err
	}
	address = l.GetAddress(newAlias)
	fmt.Fprintf(out, "==> %s [%s]\n", newAlias, address)
	return nil
}
