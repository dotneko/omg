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

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm [alias]",
	Short: "Delete an address its alias",
	Long:  `Delege an entry from the address book based on its alias.`,
	Run: func(cmd *cobra.Command, args []string) {
		l := &omg.Accounts{}

		if err := l.Load(cfg.OmgFilename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		deleted := false
		for k, a := range *l {
			if args[0] == a.Alias {
				l.DeleteIndex(k)
				if err := l.Save(cfg.OmgFilename); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				fmt.Printf("Deleted: %q [%q]\n", a.Alias, a.Address)
				deleted = true
			}
		}
		if !deleted {
			fmt.Printf("%q not found.", args[0])
		}
	},
}

func init() {
	addrCmd.AddCommand(rmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
