/*
Copyright © 2022 dotneko <EMAIL ADDRESS>

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
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
