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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [alias] [address]",
	Short: "Adds an address to the address book",
	Long: `Adds an address to the address book. For example:

Adding a normal address:
omg add nbuser onomy12345678901234567890123456789

Adding a validator address:
omg add validator onomyvaloper12345678901234567890123456789

If no alias and address provided, the an input prompt will ask
for the alias and address.
`,
	Run: func(cmd *cobra.Command, args []string) {
		l := &omg.Accounts{}

		// Read from saved address book
		if err := l.Load(cfg.OmgFilename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		alias, address, err := omg.GetAliasAddress(os.Stdin, args...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		err = l.Add(alias, address)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Save the new list
		if err := l.Save(cfg.OmgFilename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Added %q, %q to wallets\n", alias, address)
	},
}

func init() {
	addrCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
