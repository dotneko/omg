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
		amount      string
		targetDenom string
	)
	if len(args) == 2 {
		if args[1] == cfg.BaseDenom || args[1] == cfg.Token {
			amount = args[0] + args[1]
		} else {
			return fmt.Errorf("%s is not a recognized denom", args[1])
		}
	} else {
		amount = args[0]
	}
	numstr, denom, err := omg.StrSplitAmountDenom(amount)
	if err != nil {
		return err
	}
	// fmt.Fprintf(out, "Got %s, with denom %s", args[0], denom)
	if strings.EqualFold(denom, cfg.BaseDenom) {
		targetDenom = cfg.Token
	} else if strings.EqualFold(denom, cfg.Token) {
		targetDenom = cfg.BaseDenom
	} else {
		return fmt.Errorf("%s not recognized", denom)
	}

	// fmt.Fprintf(out, "Target denom: %s", targetDenom)
	convAmount := omg.ConvertAmt(numstr+denom, targetDenom)
	if detail {
		fmt.Fprintf(out, "%s => [%s] %s\n", amount, targetDenom, convAmount)
	} else {
		fmt.Fprintf(out, "%s", convAmount)
	}
	return nil
}
