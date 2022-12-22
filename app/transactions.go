package app

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	cfg "github.com/dotneko/omg/config"
	"github.com/shopspring/decimal"
)

func extractValue(line, separator, keystring string) string {
	strs := strings.Split(line, separator)
	if len(strs) == 2 {
		if strs[0] == keystring {
			return strings.TrimSpace(strs[1])
		}
	}
	return ""
}

func txGeneric(out io.Writer, cmdStr string, auto bool, keyring, outType string) (string, error) {
	if auto {
		cmdStr += " -y"
		if outType != HASH {
			fmt.Printf("Executing: %s %s\n", cfg.Daemon, cmdStr)
		}
		// Auto confirm transaction
		output, err := exec.Command(cfg.Daemon, strings.Split(cmdStr, " ")...).Output()
		if err != nil {
			return "", err
		}
		bytesReader := bytes.NewReader(output)
		bufReader := bufio.NewReader(bytesReader)
		for {
			line, _, err := bufReader.ReadLine()
			if err != nil {
				break
			}
			if outType == HASH {
				txhash := extractValue(string(line), ":", "txhash")
				if txhash != "" {
					return txhash, nil
				}
			} else {
				fmt.Fprintln(out, string(line))
			}
		}
	} else {
		// Interactive execution
		fmt.Printf("Executing: %s %s\n", cfg.Daemon, cmdStr)
		cmd := exec.Command(cfg.Daemon, strings.Split(cmdStr, " ")...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", err
		}
	}
	return "",
		nil
}

// Delegate to validator method
func TxDelegateToValidator(out io.Writer, delegator string, valAddress string, amount decimal.Decimal, auto bool, keyring, outType string) (string, error) {

	cmdStr := fmt.Sprintf("tx staking delegate %s %s --from %s", valAddress, DenomToStr(amount), delegator)
	cmdStr += fmt.Sprintf(" --fees %s --gas auto --gas-adjustment %f", cfg.DefaultFee, cfg.GasAdjust)
	cmdStr += fmt.Sprintf(" --keyring-backend %s --chain-id %s", keyring, cfg.ChainId)

	txhash, err := txGeneric(out, cmdStr, auto, keyring, outType)
	if err != nil {
		return "", err
	}
	return txhash, nil
}

// Send tokens between accounts method
func TxSend(out io.Writer, fromAddress string, toAddress string, amount decimal.Decimal, auto bool, keyring, outType string) (string, error) {

	cmdStr := fmt.Sprintf("tx bank send %s %s %s", fromAddress, toAddress, DenomToStr(amount))
	cmdStr += fmt.Sprintf(" --fees %s --gas auto --gas-adjustment %f", cfg.DefaultFee, cfg.GasAdjust)
	cmdStr += fmt.Sprintf(" --keyring-backend %s --chain-id %s", keyring, cfg.ChainId)

	txhash, err := txGeneric(out, cmdStr, auto, keyring, outType)
	if err != nil {
		return "", err
	}
	return txhash, nil
}

// Withdraw all rewards method
func TxWithdrawRewards(out io.Writer, name string, auto bool, keyring, outType string) (string, error) {

	cmdStr := fmt.Sprintf("tx distribution withdraw-all-rewards --from %s", name)
	cmdStr += fmt.Sprintf(" --fees %s --gas auto --gas-adjustment %f", cfg.DefaultFee, cfg.GasAdjust)
	cmdStr += fmt.Sprintf(" --keyring-backend %s --chain-id %s", keyring, cfg.ChainId)

	txhash, err := txGeneric(out, cmdStr, auto, keyring, outType)
	if err != nil {
		return "", err
	}
	return txhash, nil
}

// Withdraw all rewards method
func TxWithdrawValidatorCommission(out io.Writer, name string, valoperAddress string, auto bool, keyring, outType string) (string, error) {

	cmdStr := fmt.Sprintf("tx distribution withdraw-rewards %s --from %s --commission", valoperAddress, name)
	cmdStr += fmt.Sprintf(" --fees %s --gas auto --gas-adjustment %f", cfg.DefaultFee, cfg.GasAdjust)
	cmdStr += fmt.Sprintf(" --keyring-backend %s --chain-id %s", keyring, cfg.ChainId)

	txhash, err := txGeneric(out, cmdStr, auto, keyring, outType)
	if err != nil {
		return "", err
	}
	return txhash, nil
}
