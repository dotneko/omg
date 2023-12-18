/*
Copyright © 2022 dotneko
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

// delegateCmd represents the delegate command
var delegationCmd = &cobra.Command{
	Aliases: []string{"dlg", "d"},
	Use:     "delegation [account] [moniker|valoper-address]",
	Short:   "Query bonded delegation amount to validator",
	Long: `Query bonded delegation amount to validator.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		if err := cobra.ExactArgs(2)(cmd, args); err != nil {
			return fmt.Errorf("expecting [account] [moniker|valoper-address] as arguments")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		detail, err := cmd.Flags().GetBool("detail")
		if err != nil {
			return err
		}
		raw, err := cmd.Flags().GetBool("raw")
		if err != nil {
			return err
		}
		shares, err := cmd.Flags().GetBool("shares")
		if err != nil {
			return err
		}
		token, err := cmd.Flags().GetBool("token")
		if err != nil {
			return err
		}
		var outType string = ""
		if detail {
			outType = omg.DETAIL
		} else if raw {
			outType = omg.RAW
		} else if shares {
			outType = omg.SHARES
		} else if token {
			outType = omg.TOKEN
		}
		return qDelegationAction(os.Stdout, outType, args)
	},
}

func init() {
	rootCmd.AddCommand(delegationCmd)
	delegationCmd.Flags().BoolP("detail", "d", false, "Show bonded amount and delegator shares")
	delegationCmd.Flags().BoolP("raw", "r", false, "Show raw (base denom) amount")
	delegationCmd.Flags().BoolP("shares", "s", false, "Show delegator shares")
	delegationCmd.Flags().BoolP("token", "t", false, "Show token amount")

}

func qDelegationAction(out io.Writer, outType string, args []string) error {

	delegator := args[0]
	validator := args[1]
	var (
		delegatorAddress string
		valAddress       string
		valoperMoniker   string
	)

	l := &omg.Accounts{}
	if err := l.Load(cfg.OmgFilepath); err != nil {
		return err
	}
	// Check if delegator in list and is not validator account
	delegatorAddress = l.GetAddress(delegator)
	if delegatorAddress == "" {
		return fmt.Errorf("account %q not found", delegator)
	}
	if !omg.IsNormalAddress(delegatorAddress) {
		return fmt.Errorf("invalid delegator address for %s", delegator)
	}

	// Check if valid validator or validator address or moniker
	if omg.IsValidatorAddress(validator) {
		valAddress = validator
	} else {
		valAddress = l.GetAddress(validator)
		if !omg.IsValidatorAddress(valAddress) {
			// Query chain for address matching moniker if not found in address book
			searchMoniker := strings.ToLower(validator)
			valoperMoniker, valAddress = omg.GetValidator(searchMoniker)
			if valoperMoniker == "" {
				return fmt.Errorf("no validator matching %s found", validator)
			}
		}
	}
	// Check balance
	amtCoin, shares, err := omg.GetDelegationAmountShares(delegatorAddress, valAddress)
	if err != nil {
		return err
	}

	switch {
	case outType == omg.DETAIL:
		fmt.Fprintf(out, "Amount : %s ( %s%s )\n", omg.AmtToTokenStr(amtCoin.String()), omg.PrettifyBaseAmt(amtCoin.String()), cfg.BaseDenom)
		fmt.Fprintf(out, "Shares : %s\n", shares)
	case outType == omg.RAW:
		fmt.Fprintf(out, "%s\n", amtCoin.String())
	case outType == omg.SHARES:
		fmt.Fprintf(out, "Shares : %s\n", shares)
	case outType == omg.TOKEN:
		fmt.Fprintf(out, "%s\n", omg.AmtToTokenStr(amtCoin.String()))
	default:
		fmt.Fprintf(out, "Amount : %s ( %s )\n", omg.AmtToTokenStr(amtCoin.String()), omg.PrettifyBaseAmt(amtCoin.String()))
	}
	return nil
}
