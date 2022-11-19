/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// addrCmd represents the addr command
var addrCmd = &cobra.Command{
	Aliases: []string{"ad", "a"},
	Use:     "addr [command]",
	Short:   "Manage the address book",
	Long:    `Manage the address book`,
}

func init() {
	rootCmd.AddCommand(addrCmd)
}
