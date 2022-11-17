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

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename [alias] [new alias]",
	Short: "Rename an alias",
	Long: `Rename an alias. For example:

omg rename nbuser newuser
	`,
	Run: func(cmd *cobra.Command, args []string) {
		l := &omg.Accounts{}

		// Read from saved address book
		if err := l.Load(cfg.OmgFilename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if len(*l) == 0 {
			fmt.Println("No accounts in store")
			os.Exit(1)
		}
		oldAlias := args[0]
		newAlias := args[1]

		if len(newAlias) < cfg.MinAliasLength {
			fmt.Fprintln(os.Stderr, "Error: Please use alias of at least 3 characters")
			os.Exit(1)
		}
		idx := l.GetIndex(oldAlias)
		err := l.Modify(idx, newAlias, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
			os.Exit(0)
		}
		err = l.Save(cfg.OmgFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Renamed: %q to %q\n", oldAlias, newAlias)
	},
}

func init() {
	addrCmd.AddCommand(renameCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// renameCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// renameCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
