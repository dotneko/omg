package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	cfg "github.com/dotneko/omg/config"
	"github.com/dotneko/omg/types"
	"github.com/shopspring/decimal"
)

// daemon flags
const (
	jsonFlag    string = "--output json"
	keyringFlag string = "--keyring-backend"
)

// Convert denom to token (Decimal)
func DenomToTokenDec(amt decimal.Decimal) decimal.Decimal {
	d := amt.Shift(-cfg.Decimals)
	return d
}

// Convert Token to denom (Decimal)
func TokenToDenomDec(amt decimal.Decimal) decimal.Decimal {
	d := amt.Shift(cfg.Decimals)
	return d
}

func ConvertDecDenom(amount decimal.Decimal, denom string) (decimal.Decimal, string) {
	var (
		convAmount decimal.Decimal
		convDenom  string
	)
	if denom == cfg.BaseDenom {
		convAmount = DenomToTokenDec(amount)
		convDenom = cfg.Token
	} else if denom == cfg.Token {
		convAmount = TokenToDenomDec(amount)
		convDenom = cfg.BaseDenom
	}
	return convAmount, convDenom
}

// Split denominated amount to amount and denom (Decimal)
func StrSplitAmountDenomDec(amtstr string) (decimal.Decimal, string, error) {
	var NumericRegex = regexp.MustCompile(`[^0-9.-]+`)
	var AlphaRegex = regexp.MustCompile(`[^a-zA-z]+`)
	amtstr = strings.ReplaceAll(amtstr, "_", "")
	numstr := NumericRegex.ReplaceAllString(amtstr, "")
	amt, err := decimal.NewFromString(numstr)
	if err != nil {
		return decimal.NewFromInt(0), "", err
	}
	denom := AlphaRegex.ReplaceAllString(amtstr, "")
	if denom != cfg.BaseDenom && denom != cfg.Token {
		return amt, "", nil
	}
	return amt, denom, nil
}

// Convert denom to annotated string
func DenomToStr(amt decimal.Decimal) string {
	return fmt.Sprintf("%s%s", amt.String(), cfg.BaseDenom)
}

// Strip non-numeric characters and convert to decimal
func StrToDec(amtstr string) (decimal.Decimal, error) {
	var NumericRegex = regexp.MustCompile(`[^0-9.]+`)
	numstr := NumericRegex.ReplaceAllString(amtstr, "")
	amt, err := decimal.NewFromString(numstr)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	return amt, nil
}

// Insert separator for non-decimal numbers as output
func PrettifyDenom(amt decimal.Decimal) string {
	if amt.LessThan(decimal.NewFromInt(1000)) {
		return amt.String()
	}
	var (
		amtStr string
		//decimalStr string
		dotPos int
	)
	dotPos = strings.Index(amt.String(), ".")
	if dotPos == -1 {
		amtStr = amt.String()
		// decimalStr = ""
	} else {
		s := strings.Split(amt.String(), ".")
		amtStr = s[0]
		// decimalStr = s[1]
	}
	separator := "_"
	startIdx := len(amtStr) % 3
	if startIdx == 0 {
		startIdx = 3
	}
	outStr := amtStr[:startIdx]
	pos := startIdx
	for pos < len(amtStr) {
		outStr = outStr + separator + amtStr[pos:pos+3]
		pos = pos + 3
	}
	if dotPos == -1 {
		return outStr
	}
	return outStr + "._"
}

func PrettifyAmount(amount decimal.Decimal, denom string) string {
	if denom == cfg.BaseDenom {
		return fmt.Sprintf("%s %s", PrettifyDenom(amount), denom)
	}
	if denom == cfg.Token {
		return fmt.Sprintf("%s %s", amount.String(), denom)
	}
	return ""
}

// Get Balances Query
func GetBalancesQuery(address string) (*types.BalancesQuery, error) {
	cmdStr := fmt.Sprintf("query bank balances %s %s", jsonFlag, address)
	out, err := exec.Command(cfg.Daemon, strings.Split(cmdStr, " ")...).Output()
	if err != nil {
		return nil, fmt.Errorf("cannot get balance for %s", address)
	}
	if !json.Valid(out) {
		return nil, errors.New("invalid json")
	}
	var b types.BalancesQuery
	if err = json.Unmarshal(out, &b); err != nil {
		return nil, err
	}

	return &b, nil
}

// Get Balances (first denom) to decimal amount
func GetBalanceDec(address string) (decimal.Decimal, error) {

	bQ, err := GetBalancesQuery(address)
	if err != nil {
		return decimal.NewFromInt(-1), err
	}
	amt, err := decimal.NewFromString(bQ.Balances[0].Amount)
	if err != nil {
		return decimal.NewFromInt(-1), err
	}
	return amt, nil
}

// Check balance method
func CheckBalances(address string) {
	balance, err := GetBalanceDec(address)
	if err != nil {
		fmt.Sprintln(err)
	}
	fmt.Printf("Avaliable balance : %s %s (%s %s)\n", balance.String(), cfg.BaseDenom, DenomToTokenDec(balance).String(), cfg.Token)
}

// Get keyring name and addresses
func GetKeyringAccounts(keyring string) (Accounts, error) {

	cmdStr := fmt.Sprintf("keys list %s %s %s", jsonFlag, keyringFlag, keyring)
	out, err := exec.Command(cfg.Daemon, strings.Split(cmdStr, " ")...).Output()
	if err != nil {
		return nil, err
	}
	if !json.Valid(out) {
		return nil, errors.New("invalid json")
	}
	var k []types.KeyStruct
	if err = json.Unmarshal(out, &k); err != nil {
		return nil, err
	}
	if len(k) == 0 {
		return nil, errors.New("no addresses in keyring")
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
		return nil, fmt.Errorf("invalid json")
	}
	var r types.RewardsQuery
	if err = json.Unmarshal(out, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

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
	cmdStr += fmt.Sprintf(" --fees %s--gas auto --gas-adjustment %f", cfg.DefaultFee, cfg.GasAdjust)
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
