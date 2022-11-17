/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// addrCmd represents the addr command
var addrCmd = &cobra.Command{
	Use:   "addr [command]",
	Short: "Manage the address book",
	Long: `Manage the address book
	
Add addresses with command: add [alias] [address]
Delete addresses with command: rm [alias]
Rename alias with command: rename [alias] [new alias]
List alias and address with command: list`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("addr called")
	// },
}

func init() {
	rootCmd.AddCommand(addrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
