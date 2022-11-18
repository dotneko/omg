/*
Copyright © 2022 dotneko <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"

	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var (
	convertCmd = &cobra.Command{
		Use:   "convert [amount][denom]",
		Short: fmt.Sprintf("Conversion between %s and %s", cfg.Denom, cfg.Token),
		Long:  `Conversion tool between token and denom amounts`,
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
		Run: func(cmd *cobra.Command, args []string) {
			var (
				amount float64
				denom  string
				err    error
			)
			if len(args) == 1 {
				amount, denom, err = omg.StrSplitAmountDenom(args[0])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			} else if len(args) == 2 {
				amount, err = omg.StrToFloat(args[0])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				denom = args[1]
			}
			fmt.Printf("%.f %s => ", amount, denom)
			var (
				outFormat  string
				convAmount float64
				convDenom  string
			)
			if denom == cfg.Denom {
				convAmount = omg.DenomToToken(amount)
				convDenom = cfg.Token
				outFormat = "%.18f%s"
			} else if denom == cfg.Token {
				convAmount = omg.TokenToDenom(amount)
				convDenom = cfg.Denom
				outFormat = "%.0f%s"
			}
			fmt.Printf(outFormat, convAmount, convDenom)

		},
	}
)

func init() {
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
