/*
Copyright © 2022 dotneko <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"
	"github.com/shopspring/decimal"

	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var (
	convertCmd = &cobra.Command{
		Aliases: []string{"cv", "c"},
		Use:     "convert [amount][denom]",
		Short:   fmt.Sprintf("Conversion between %s and %s", cfg.BaseDenom, cfg.Token),
		Long:    fmt.Sprintf("Conversion between %s and %s", cfg.BaseDenom, cfg.Token),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
			if err := cobra.RangeArgs(1, 2)(cmd, args); err != nil {
				return fmt.Errorf("expecting [amount][denom]")
			}
			if len(args) == 2 {
				denomStr := strings.ToLower(args[1])
				if denomStr != cfg.BaseDenom && denomStr != strings.ToLower(cfg.Token) {
					return fmt.Errorf("denom must be %q or %q", cfg.BaseDenom, cfg.Token)
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			detail, err := cmd.Flags().GetBool("detail")
			if err != nil {
				return err
			}
			return convertAction(os.Stdout, detail, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().BoolP("detail", "d", false, "Output detail")

}

func convertAction(out io.Writer, detail bool, args []string) error {
	var (
		amount decimal.Decimal
		denom  string
		err    error
	)
	if len(args) == 1 {
		amount, denom, err = omg.StrSplitAmountDenomDec(args[0])
		if err != nil {
			return err
		}
	} else if len(args) == 2 {
		amount, err = decimal.NewFromString(args[0])
		if err != nil {
			return err
		}
		denom = args[1]
	}
	var (
		convAmount decimal.Decimal
		convDenom  string
	)
	convAmount, convDenom = omg.ConvertDecDenom(amount, denom)
	if detail {
		fmt.Fprintf(out, "%s => %s\n", omg.PrettifyAmount(amount, denom), omg.PrettifyAmount(convAmount, convDenom))
	} else {
		fmt.Fprintf(out, "%s%s ", convAmount.String(), convDenom)
	}
	return nil
}
