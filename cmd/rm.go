/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"

	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Aliases: []string{"remove", "delete", "r"},
	Use:     "rm [name]",
	Short:   "Delete an address book entry",
	Long: `Delete an address book entry.
	
Examples:

To delete a user (with confirmation):
# omg address rm user1

To delete a user (force confirmation):
# omg address rm user1 -f
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("expecting [name] as argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}
		return rmAction(os.Stdin, os.Stdout, force, args)
	},
}

func init() {
	addressCmd.AddCommand(rmCmd)
	rmCmd.Flags().BoolP("force", "f", false, "Force delete without confirmation")

}

func rmAction(in io.Reader, out io.Writer, force bool, args []string) error {
	l := &omg.Accounts{}

	if err := l.Load(cfg.OmgFilepath); err != nil {
		return err
	}

	idx, address := l.GetIndexAddress(args[0])
	if idx == -1 {
		return fmt.Errorf("%q not found in address book", args[0])
	}
	fmt.Fprintf(out, "Found: %s [%s]\n", args[0], address)
	confirm := ""
	if !force {
		s := bufio.NewScanner(in)
		fmt.Fprintf(out, "Confirm delete [Y/n]?")
		s.Scan()
		confirm = strings.ToLower(s.Text())
	}
	if confirm == "" || confirm == "y" {
		l.DeleteIndex(idx)
		if err := l.Save(cfg.OmgFilepath); err != nil {
			return err
		}
		fmt.Fprintf(out, "Deleted: %q [%s]\n", args[0], address)
	}
	return nil
}
