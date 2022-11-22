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
	Aliases: []string{"mod", "mv", "m", "rename"},
	Use:     "modify [name] [new name] *OPTIONAL -a [new address]",
	Short:   "Modify an address book entry",
	Long: `Modify an address book entry
	
Examples:

# omg addr modify user1 user2
# omg addr modify user1 user2 -a newaddress1234567890123456789
# omg addr modify user1 -a newaddress1234567890123456789`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		if err := cobra.RangeArgs(1, 2)(cmd, args); err != nil {
			return fmt.Errorf("expecting [name] [new name] as arguments")
		}
		if omg.IsValidAddress(args[0]) {
			return fmt.Errorf("[name] cannot be an address")
		}
		if len(args) > 1 && omg.IsValidAddress(args[1]) {
			return fmt.Errorf("[new name] cannot be an address")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("address")
		if err != nil {
			return err
		}
		if len(args) == 1 && address == "" {
			return fmt.Errorf("missing new name or address flag")
		}
		return modifyAction(os.Stdout, address, args)
	},
}

func init() {
	addressCmd.AddCommand(modifyCmd)

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
		return fmt.Errorf("no accounts in store")
	}

	if len(newAlias) < cfg.MinAliasLength {
		return fmt.Errorf("please use alias of at least 3 characters")
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
