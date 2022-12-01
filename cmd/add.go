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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Aliases: []string{"a"},
	Use:     "add [name] [address]",
	Short:   "Add an address to the address book",
	Long: fmt.Sprintf(`Adds an entry to the address book

The entry [name] may be an alias for external addresses, or match the keyring-backend name
for user-owned accounts.

For transactions to process, the [name] is used as the parameter by the %s daemon
to check the keyring for available user-owned accounts to generate signature for signing.

The entry [address] must not be a valoper address.

`, cfg.Daemon),
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		if err := cobra.ExactArgs(2)(cmd, args); err != nil {
			return fmt.Errorf("expecting [name] [address] as arguments")
		}
		if omg.IsValidAddress(args[0]) {
			return fmt.Errorf("[name] cannot be an address")
		}
		if !omg.IsNormalAddress(args[1]) {
			return fmt.Errorf("%s is not a valid address", args[1])
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return addAction(os.Stdout, args)
	},
}

func init() {
	addressCmd.AddCommand(addCmd)
}

func addAction(out io.Writer, args []string) error {
	l := &omg.Accounts{}

	// Read from saved address book
	if err := l.Load(cfg.OmgFilepath); err != nil {
		return err
	}
	alias := args[0]
	address := args[1]

	err := l.Add(alias, address)
	if err != nil {
		return err
	}
	// Save the new list
	if err := l.Save(cfg.OmgFilepath); err != nil {
		return err
	}
	fmt.Fprintf(out, "Added ==> %s [%s]\n", alias, address)
	return nil
}
