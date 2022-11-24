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

// validatorQueryCmd represents the validatorQuery command
var validatorQueryCmd = &cobra.Command{
	Aliases: []string{"show", "q", "s"},
	Use:     "query",
	Short:   "Query staking validators",
	Long:    `Query staking validators.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			cmd.Help()
			os.Exit(0)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return validatorQueryAction(os.Stdout, args)
	},
}

func init() {
	validatorCmd.AddCommand(validatorQueryCmd)
}

func validatorQueryAction(out io.Writer, args []string) error {
	vQ, err := omg.GetValidatorsQuery()
	if err != nil {
		return err
	}
	for _, val := range vQ.Validators {
		if !val.Jailed {
			fmt.Fprintf(out, "%20s : %s\n", val.Description.Moniker, val.OperatorAddress)
		}
	}
	l := &omg.Accounts{}

	if err := l.Load(cfg.OmgFilepath); err != nil {
		return err
	}
	return nil
}
