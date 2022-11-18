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

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm [alias]",
	Short: "Delete an address its alias",
	Long:  `Delege an entry from the address book based on its alias.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rmAction(os.Stdout, args)
	},
}

func init() {
	addrCmd.AddCommand(rmCmd)
}

func rmAction(out io.Writer, args []string) error {
	l := &omg.Accounts{}

	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	deleted := false
	for k, a := range *l {
		if args[0] == a.Alias {
			l.DeleteIndex(k)
			if err := l.Save(cfg.OmgFilename); err != nil {
				return err
			}
			fmt.Fprintf(out, "Deleted: %q [%q]\n", a.Alias, a.Address)
			deleted = true
		}
	}
	if !deleted {
		return fmt.Errorf("%q not found.", args[0])
	}
}
