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

The entry [address] can be a normal address or a validator/valoper address.

Examples:

Adding a normal address:
# omg addr add nbuser %s12345678901234567890123456789

Adding a validator/valoper address:
# omg addr add validator %s12345678901234567890123456789

An input prompt would ask for the [name] and [address] if these are not specified.
`, cfg.Daemon, cfg.AddressPrefix, cfg.ValoperPrefix),
	RunE: func(cmd *cobra.Command, args []string) error {
		return addAction(os.Stdout, args)
	},
}

func init() {
	addrCmd.AddCommand(addCmd)
}

func addAction(out io.Writer, args []string) error {
	l := &omg.Accounts{}

	// Read from saved address book
	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	alias, address, err := omg.GetAliasAddress(os.Stdin, args...)
	if err != nil {
		return err
	}
	err = l.Add(alias, address)
	if err != nil {
		return err
	}
	// Save the new list
	if err := l.Save(cfg.OmgFilename); err != nil {
		return err
	}
	fmt.Fprintf(out, "Added ==> %s [%s]\n", alias, address)
	return nil
}
