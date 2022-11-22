/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"io"
	"os"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Aliases: []string{"imp", "i"},
	Use:     "import",
	Short:   "Import addresses from keyring",
	Long:    `Import addresses from keyring`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			cmd.Help()
			os.Exit(0)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		keyring, err := cmd.Flags().GetString("keyring")
		if err != nil {
			return err
		}
		return importAction(os.Stdout, keyring, args)
	},
}

func init() {
	addressCmd.AddCommand(importCmd)

	importCmd.Flags().StringP("keyring", "k", cfg.KeyringBackend, "Specify keyring-backend")
}

func importAction(out io.Writer, keyring string, args []string) error {
	l := &omg.Accounts{}
	if err := l.Load(cfg.OmgFilename); err != nil {
		return err
	}
	num, err := omg.ImportFromKeyring(l, keyring)
	if err != nil {
		return err
	}
	if num > 0 {
		// Save the new list
		if err := l.Save(cfg.OmgFilename); err != nil {
			return err
		}
	}
	return nil
}
