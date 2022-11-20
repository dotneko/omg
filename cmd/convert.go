/*
Copyright © 2022 dotneko <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io"
	"os"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"
	"github.com/shopspring/decimal"

	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var (
	convertCmd = &cobra.Command{
		Aliases: []string{"cv"},
		Use:     "convert [amount][denom]",
		Short:   fmt.Sprintf("Conversion between %s and %s", cfg.Denom, cfg.Token),
		Long:    `Conversion tool between token and denom amounts`,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.RangeArgs(1, 2)(cmd, args); err != nil {
				return err
			}
			if len(args) == 2 {
				if args[1] != cfg.Denom && args[1] != cfg.Token {
					fmt.Printf("Error: denom must be %q or %q; got %q\n", cfg.Denom, cfg.Token, args[1])
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return convertAction(os.Stdout, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(convertCmd)
}

func convertAction(out io.Writer, args []string) error {
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
	if denom == cfg.Denom {
		convAmount = omg.DenomToTokenDec(amount)
		convDenom = cfg.Token
	} else if denom == cfg.Token {
		convAmount = omg.TokenToDenomDec(amount)
		convDenom = cfg.Denom
	}
	fmt.Fprintf(out, "%s%s", convAmount.String(), convDenom)
	return nil
}
