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

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Aliases: []string{"ren", "mv"},
	Use:     "rename [alias] [new alias]",
	Short:   "Rename an alias",
	Long:    `Rename an alias. For example:`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("address")
		if err != nil {
			return err
		}
		return renameAction(os.Stdout, address, args)
	},
}

func init() {
	addrCmd.AddCommand(renameCmd)

	renameCmd.Flags().StringP("address", "a", "", "Address to modify to")
}

func renameAction(out io.Writer, address string, args []string) error {
	l := &omg.Accounts{}

	// Read from saved address book
	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	if len(*l) == 0 {
		return fmt.Errorf("No accounts in store")
	}
	oldAlias := args[0]
	newAlias := args[1]

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
	fmt.Fprintf(out, "Renamed: %q to %q\n", oldAlias, newAlias)
	return nil
}
