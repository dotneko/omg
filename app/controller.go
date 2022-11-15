package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	cfg "github.com/dotneko/omg/config"
	"github.com/dotneko/omg/types"
)

// daemon flags
const (
	jsonFlag    string = "--output json"
	keyringFlag string = "--keyring-backend"
)

// Convert denom to token amount
func DenomToToken(amt float64) float64 {
	return amt / math.Pow10(cfg.Decimals)
}

// Convert denom to annotated string
func DenomToStr(amt float64) string {
	return fmt.Sprintf("%.0f%s", amt, cfg.Denom)
}

// Convert token amount to denom amount
func TokenToDenom(amt float64) float64 {
	return amt * math.Pow10(cfg.Decimals)
}

// Strip non-numeric characters and convert to float
func StrToFloat(amtstr string) (float64, error) {
	var NumericRegex = regexp.MustCompile(`[^0-9.]+`)
	numstr := NumericRegex.ReplaceAllString(amtstr, "")
	amt, err := strconv.ParseFloat(numstr, 64)
	if err != nil {
		return 0, err
	}
	return amt, nil
}

// Split denominated amount to amount and denom
func StrSplitAmountDenom(amtstr string) (float64, string, error) {
	var NumericRegex = regexp.MustCompile(`[^0-9.-]+`)
	var AlphaRegex = regexp.MustCompile(`[^a-zA-z]+`)
	numstr := NumericRegex.ReplaceAllString(amtstr, "")
	amt, err := strconv.ParseFloat(numstr, 64)
	if err != nil {
		return 0, "", err
	}
	denom := AlphaRegex.ReplaceAllString(amtstr, "")
	if denom != cfg.Denom && denom != cfg.Token {
		return amt, "", nil
	}
	return amt, denom, nil
}

// Get Balances Query
func GetBalancesQuery(address string) (*types.BalancesQuery, error) {
	cmdStr := fmt.Sprintf("query bank balances %s %s", jsonFlag, address)
	out, err := exec.Command(cfg.Daemon, strings.Split(cmdStr, " ")...).Output()
	if err != nil {
		return nil, err
	}
	if !json.Valid(out) {
		return nil, errors.New("Invalid json")
	}
	var b types.BalancesQuery
	if err = json.Unmarshal(out, &b); err != nil {
		return nil, err
	}

	return &b, nil
}

// Get Balances (first denom) to float amount
func GetBalanceAmount(address string) (float64, error) {

	bQ, err := GetBalancesQuery(address)
	if err != nil {
		return -1, err
	}
	amt, err := strconv.ParseFloat(bQ.Balances[0].Amount, 64)
	if err != nil {
		return -1, err
	}
	return amt, nil
}

// Get keyring name and addresses
func GetKeyringAccounts(keyring string) (Accounts, error) {

	cmdStr := fmt.Sprintf("keys list %s %s %s", jsonFlag, keyringFlag, keyring)
	out, err := exec.Command(cfg.Daemon, strings.Split(cmdStr, " ")...).Output()
	if err != nil {
		return nil, err
	}
	if !json.Valid(out) {
		return nil, errors.New("Invalid json")
	}
	var k []types.KeyStruct
	if err = json.Unmarshal(out, &k); err != nil {
		return nil, err
	}
	if len(k) == 0 {
		return nil, errors.New("No addresses in keyring")
	}
	var accounts []Account = nil
	for _, key := range k {
		acc := Account{
			Alias:   key.Name,
			Address: key.Address,
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

// Parse rewards
func GetRewards(address string) (*types.RewardsQuery, error) {

	cmdStr := fmt.Sprintf("query distribution rewards %s %s", jsonFlag, address)
	out, err := exec.Command(cfg.Daemon, strings.Split(cmdStr, " ")...).Output()
	if err != nil {
		return nil, err
	}
	if !json.Valid(out) {
		return nil, fmt.Errorf("Invalid json")
	}
	var r types.RewardsQuery
	if err = json.Unmarshal(out, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Delegate to validator method
func TxDelegateToValidator(delegator string, valAddress string, amount float64, keyring string, auto bool) error {

	cmdStr := fmt.Sprintf("tx staking delegate %s %s --from %s", valAddress, DenomToStr(amount), delegator)
	cmdStr += fmt.Sprintf(" --fees %d%s --gas auto --gas-adjustment %f", cfg.DefaultFee, cfg.Denom, cfg.GasAdjust)
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
func TxSend(fromAddress string, toAddress string, amount float64, keyring string, auto bool) error {
	// fmt.Printf("DelegateToValidator %s %s %s %t\n", delegator, valAddress, denomToStr(amount), auto)

	cmdStr := fmt.Sprintf("tx bank send %s %s %s", fromAddress, toAddress, DenomToStr(amount))
	//cmdStr += fmt.Sprintf(" --fees %d%s --gas auto --gas-adjustment %f", defaultFee, denom, gasAdjust)
	cmdStr += fmt.Sprintf("--gas auto --gas-adjustment %f", cfg.GasAdjust)
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
func TxWithdrawRewards(alias string, keyring string, auto bool) error {

	cmdStr := fmt.Sprintf("tx distribution withdraw-all-rewards --from %s", alias)
	cmdStr += fmt.Sprintf(" --fees %d%s --gas auto --gas-adjustment %f", cfg.DefaultFee, cfg.Denom, cfg.GasAdjust)
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
