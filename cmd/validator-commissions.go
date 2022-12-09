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

// commissionsCmd represents the commissions command
var commissionsCmd = &cobra.Command{
	Aliases: []string{"com", "c"},
	Use:     "commissions [moniker|valoper-address]",
	Short:   "Query validator commission",
	Long:    `Query validator commission.`,
	Args: func(cmd *cobra.Command, args []string) error {

		if err := cobra.RangeArgs(0, 1)(cmd, args); err != nil {
			return fmt.Errorf("expecting [moniker|valoper-address] as argument or --all flag")
		}
		allAccounts, _ := cmd.Flags().GetBool("all")
		if len(args) == 0 && !allAccounts {
			cmd.Help()
			os.Exit(0)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		allAccounts, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}
		detail, err := cmd.Flags().GetBool("detail")
		if err != nil {
			return err
		}
		raw, err := cmd.Flags().GetBool("raw")
		if err != nil {
			return err
		}
		token, err := cmd.Flags().GetBool("token")
		if err != nil {
			return err
		}
		return commissionsAction(os.Stdout, allAccounts, detail, raw, token, args)
	},
}

func init() {
	validatorCmd.AddCommand(commissionsCmd)

	commissionsCmd.Flags().BoolP("all", "a", false, "Check all validators")
	commissionsCmd.Flags().BoolP("detail", "d", false, "Detailed output")
	commissionsCmd.Flags().BoolP("raw", "r", false, "Raw output")
	commissionsCmd.Flags().BoolP("token", "t", false, "Token amount output")

}

func commissionsAction(out io.Writer, allAccounts, detail, raw, token bool, args []string) error {
	var (
		moniker        string = ""
		valoperAddress string = ""
		outType        string = ""
	)
	if raw {
		outType = omg.RAW
	} else if token {
		outType = omg.TOKEN
	} else if detail {
		outType = omg.DETAIL
	}
	if len(args) == 1 {
		moniker, valoperAddress = omg.GetValidator(args[0])

		if valoperAddress == "" {
			return fmt.Errorf("no matching validator for %s", args[0])
		}

		commission, err := omg.GetCommissionDec(valoperAddress)
		if err != nil {
			return err
		}
		omg.OutputAmount(out, moniker, valoperAddress, commission, cfg.BaseDenom, outType)
		return nil
	}

	if allAccounts {
		vQ, err := omg.GetValidatorsQuery()
		if err != nil {
			return err
		}
		for _, val := range vQ.Validators {
			if !val.Jailed {
				commission, err := omg.GetCommissionDec(val.OperatorAddress)
				if err != nil {
					fmt.Fprintf(out, "Error: %12s : %s\n", val.Description.Moniker, err.Error())
					continue
				}
				if outType == omg.RAW || outType == omg.TOKEN {
					fmt.Fprintf(out, "%s ", val.OperatorAddress)
				}
				omg.OutputAmount(out, val.Description.Moniker, val.OperatorAddress, commission, cfg.BaseDenom, outType)
			}
		}
	}
	return nil
}
