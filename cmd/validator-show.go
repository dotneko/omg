/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io"
	"os"

	omg "github.com/dotneko/omg/app"
	"github.com/spf13/cobra"
)

// validatorShowCmd represents the validatorShow command
var validatorShowCmd = &cobra.Command{
	Aliases: []string{"list", "s"},
	Use:     "show **OPTIONAL:[moniker]",
	Short:   "Show staking validators",
	Long:    `Show staking validators.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.RangeArgs(0, 1)(cmd, args); err != nil {
			fmt.Printf("Error: %s\n", fmt.Errorf("too many arguments"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return validatorShowAction(os.Stdout, args)
	},
}

func init() {
	validatorCmd.AddCommand(validatorShowCmd)
}

func validatorShowAction(out io.Writer, args []string) error {
	if len(args) == 1 {
		_, valoperAddress := omg.GetValidator(args[0])
		if valoperAddress != "" {
			fmt.Fprintln(out, valoperAddress)
		}
		return nil
	}
	vQ, err := omg.GetValidatorsQuery()
	if err != nil {
		return err
	}
	for _, val := range vQ.Validators {
		if !val.Jailed {
			fmt.Fprintf(out, "%20s [ %s ]\n", val.Description.Moniker, val.OperatorAddress)
		}
	}
	return nil
}
