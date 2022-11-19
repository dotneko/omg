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
	Use:     "query [balances | rewards] [name]",
	Short:   "Query balances or rewards for an name/address",
	Long: `Query balances or rewards for an name/address.
	
Examples:
# omg query balances user1
# omg query rewards user1

To query all entries in address book:
# omg query balances -a
# omg query rewards -a
`,
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
