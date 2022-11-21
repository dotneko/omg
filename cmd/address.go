/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// addrCmd represents the addr command
var addressCmd = &cobra.Command{
	Aliases: []string{"addr", "ad", "a"},
	Use:     "address [command]",
	Short:   "Manage the address book",
	Long:    `Manage the address book`,
}

func init() {
	rootCmd.AddCommand(addressCmd)
}
