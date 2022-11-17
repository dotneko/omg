/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Aliases: []string{"q"},
	Use:     "query [balance/reward] [alias]",
	Short:   "Query balance or reward for an address/alias",
	Long:    `Query balance or reward for an address/alias.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("query called")
	// },
}

func init() {
	rootCmd.AddCommand(queryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// queryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// queryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
