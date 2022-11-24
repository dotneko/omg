/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"
	"github.com/spf13/cobra"
)

// validatorImportCmd represents the validatorImport command
var validatorImportCmd = &cobra.Command{
	Aliases: []string{"imp", "i"},
	Use:     "import [moniker]",
	Short:   "Import an validator to the address book and optionally assign an alias",
	Long:    `Import an validator to the address book and optionally assign an alias.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("expecting [moniker] as argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		alias, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		return validatorImportAction(os.Stdout, alias, args)
	},
}

func init() {
	validatorCmd.AddCommand(validatorImportCmd)

	validatorImportCmd.Flags().StringP("name", "n", "", "Store specified name/alias instead of moniker")
}

func validatorImportAction(out io.Writer, alias string, args []string) error {

	searchMoniker := strings.ToLower(args[0])
	valoperMoniker, valoperAddress := omg.GetValidatorAddress(searchMoniker)

	l := &omg.Accounts{}
	if err := l.Load(cfg.OmgFilepath); err != nil {
		return err
	}

	if alias == "" {
		alias = valoperMoniker
	}
	err := l.Add(alias, valoperAddress)
	if err != nil {
		return err
	}
	// Save the new list
	if err := l.Save(cfg.OmgFilepath); err != nil {
		return err
	}
	fmt.Fprintf(out, "Added ==> %s [%s]\n", alias, valoperAddress)
	return nil
}
