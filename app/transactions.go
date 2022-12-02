package app

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	cfg "github.com/dotneko/omg/config"
	"github.com/shopspring/decimal"
)

// Delegate to validator method
func TxDelegateToValidator(delegator string, valAddress string, amount decimal.Decimal, keyring string, auto bool) error {

	cmdStr := fmt.Sprintf("tx staking delegate %s %s --from %s", valAddress, DenomToStr(amount), delegator)
	cmdStr += fmt.Sprintf(" --fees %s --gas auto --gas-adjustment %f", cfg.DefaultFee, cfg.GasAdjust)
	cmdStr += fmt.Sprintf(" --keyring-backend %s --chain-id %s", keyring, cfg.ChainId)

	fmt.Printf("Executing: %s %s\n", cfg.Daemon, cmdStr)
	cmd := exec.Command(cfg.Daemon, strings.Split(cmdStr, " ")...)

	if auto {
		// Auto confirm transaction
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			return err
		}
		// Expect prompt and confirm with 'y'
		stdin.Write([]byte("y\n"))

		if err := cmd.Wait(); err != nil {
			return err
		}
	} else {
		// Interactive execution
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// Send tokens between accounts method
func TxSend(fromAddress string, toAddress string, amount decimal.Decimal, keyring string, auto bool) error {
	// fmt.Printf("DelegateToValidator %s %s %s %t\n", delegator, valAddress, denomToStr(amount), auto)

	cmdStr := fmt.Sprintf("tx bank send %s %s %s", fromAddress, toAddress, DenomToStr(amount))
	//cmdStr += fmt.Sprintf(" --fees %d%s --gas auto --gas-adjustment %f", defaultFee, denom, gasAdjust)
	cmdStr += fmt.Sprintf(" --fees %s --gas auto --gas-adjustment %f", cfg.DefaultFee, cfg.GasAdjust)
	cmdStr += fmt.Sprintf(" --keyring-backend %s --chain-id %s", keyring, cfg.ChainId)

	fmt.Printf("Executing: %s %s\n", cfg.Daemon, cmdStr)
	cmd := exec.Command(cfg.Daemon, strings.Split(cmdStr, " ")...)

	if auto {
		// Auto confirm transaction
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			return err
		}
		// Expect prompt and confirm with 'y'
		stdin.Write([]byte("y\n"))

		if err := cmd.Wait(); err != nil {
			return err
		}
	} else {
		// Interactive execution
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// Withdraw all rewards method
func TxWithdrawRewards(out io.Writer, name string, keyring string, auto bool) error {

	cmdStr := fmt.Sprintf("tx distribution withdraw-all-rewards --from %s", name)
	cmdStr += fmt.Sprintf(" --fees %s --gas auto --gas-adjustment %f", cfg.DefaultFee, cfg.GasAdjust)
	cmdStr += fmt.Sprintf(" --keyring-backend %s --chain-id %s", keyring, cfg.ChainId)

	fmt.Fprintf(out, "Executing: %s %s\n", cfg.Daemon, cmdStr)
	cmd := exec.Command(cfg.Daemon, strings.Split(cmdStr, " ")...)

	if auto {
		// Auto confirm transaction
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}

		cmd.Stdout = out
		cmd.Stderr = out
		if err := cmd.Start(); err != nil {
			return err
		}
		// Expect prompt and confirm with 'y'
		stdin.Write([]byte("y\n"))

		if err := cmd.Wait(); err != nil {
			return err
		}
	} else {
		// Interactive execution
		cmd.Stdin = os.Stdin
		cmd.Stdout = out
		cmd.Stderr = out
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// Withdraw all rewards method
func TxWithdrawValidatorCommission(out io.Writer, name string, valoperAddress string, keyring string, auto bool) error {

	cmdStr := fmt.Sprintf("tx distribution withdraw-rewards %s --from %s --commission", valoperAddress, name)
	cmdStr += fmt.Sprintf(" --fees %s --gas auto --gas-adjustment %f", cfg.DefaultFee, cfg.GasAdjust)
	cmdStr += fmt.Sprintf(" --keyring-backend %s --chain-id %s", keyring, cfg.ChainId)

	fmt.Fprintf(out, "Executing: %s %s\n", cfg.Daemon, cmdStr)
	cmd := exec.Command(cfg.Daemon, strings.Split(cmdStr, " ")...)

	if auto {
		// Auto confirm transaction
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}

		cmd.Stdout = out
		cmd.Stderr = out
		if err := cmd.Start(); err != nil {
			return err
		}
		// Expect prompt and confirm with 'y'
		stdin.Write([]byte("y\n"))

		if err := cmd.Wait(); err != nil {
			return err
		}
	} else {
		// Interactive execution
		cmd.Stdin = os.Stdin
		cmd.Stdout = out
		cmd.Stderr = out
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
