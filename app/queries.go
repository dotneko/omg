package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	cfg "github.com/dotneko/omg/config"

	"github.com/dotneko/omg/types"
)

// daemon flags
const (
	jsonFlag    string = "--output json"
	keyringFlag string = "--keyring-backend"
)

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

// Get Balances (first denom) to string amount
func GetBalance(address string) (sdktypes.Coin, error) {

	bQ, err := GetBalancesQuery(address)
	if len(bQ.Balances) == 0 {
		return sdktypes.Coin{}, fmt.Errorf("no balances found")
	}
	if err != nil {
		return sdktypes.Coin{}, err
	}
	amtCoin, err := sdktypes.ParseCoinNormalized(bQ.Balances[0].Amount + bQ.Balances[0].Denom)
	if err != nil {
		return sdktypes.Coin{}, err
	}
	return amtCoin, nil
}

// Check balance method
func CheckBalances(address string) {
	balance, err := GetBalance(address)
	if err != nil {
		fmt.Sprintln(err)
	}
	tokenBalance, err := AmtToTokenDecCoin(balance.String())
	if err != nil {
		fmt.Sprintln(err)
	}
	fmt.Printf("Avaliable balance : %s %s (%s %s)\n", tokenBalance.Amount.String(), cfg.BaseDenom, tokenBalance, cfg.Token)
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

// Query Keyring for address
func QueryKeyringAddress(name, keyring string) string {

	cmdStr := fmt.Sprintf("keys show %s %s %s %s", name, jsonFlag, keyringFlag, keyring)
	out, err := exec.Command(cfg.Daemon, strings.Split(cmdStr, " ")...).CombinedOutput()
	if err != nil {
		return ""
	}
	if strings.Contains(string(out), "Error") {
		return ""
	}
	if !json.Valid(out) {
		return ""
	}
	var k types.KeyStruct
	if err = json.Unmarshal(out, &k); err != nil {
		return ""
	}
	return k.Address
}

// Query Validators
func GetValidatorsQuery() (*types.ValidatorsQuery, error) {
	cmdStr := fmt.Sprintf("query staking validators %s", jsonFlag)
	out, err := exec.Command(cfg.Daemon, strings.Split(cmdStr, " ")...).Output()
	if err != nil {
		return nil, err
	}
	if !json.Valid(out) {
		return nil, errors.New("invalid json")
	}
	var v types.ValidatorsQuery
	if err = json.Unmarshal(out, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Search validators by moniker or valoper-address
func GetValidator(search string) (string, string) {
	searchStr := strings.ToLower(search)

	vQ, _ := GetValidatorsQuery()

	for _, val := range vQ.Validators {
		if !val.Jailed &&
			searchStr == strings.ToLower(val.Description.Moniker) ||
			searchStr == strings.ToLower(val.OperatorAddress) {
			return val.Description.Moniker, val.OperatorAddress
		}
	}
	return "", ""
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

// Parse commissions for validator
func QueryCommission(valopAddress string) (*types.CommissionsQuery, error) {

	cmdStr := fmt.Sprintf("query distribution commission %s %s", jsonFlag, valopAddress)
	out, err := exec.Command(cfg.Daemon, strings.Split(cmdStr, " ")...).Output()
	if err != nil {
		return nil, err
	}
	if !json.Valid(out) {
		return nil, fmt.Errorf("invalid json")
	}
	var c types.CommissionsQuery
	if err = json.Unmarshal(out, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// Get commissions (first denom) to decimal amount
func GetCommission(valopAddress string) (sdktypes.Coin, error) {

	cQ, err := QueryCommission(valopAddress)
	if len(cQ.Commission) == 0 {
		return sdktypes.Coin{}, fmt.Errorf("no commission found")
	}
	if err != nil {
		return sdktypes.Coin{}, err
	}
	amtCoin, err := sdktypes.ParseCoinNormalized(cQ.Commission[0].Amount + cQ.Commission[0].Denom)
	if err != nil {
		return sdktypes.Coin{}, err
	}
	return amtCoin, nil
}

// Query delegation
func QueryDelegation(address, valopAddress string) (*types.DelegationQuery, error) {

	cmdStr := fmt.Sprintf("query staking delegation %s %s %s", address, valopAddress, jsonFlag)
	out, err := exec.Command(cfg.Daemon, strings.Split(cmdStr, " ")...).Output()
	if err != nil {
		return nil, err
	}
	if !json.Valid(out) {
		return nil, fmt.Errorf("invalid json")
	}
	var d types.DelegationQuery
	if err = json.Unmarshal(out, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

// Get Delegation Amount
func GetDelegationAmountShares(address, valopAddress string) (sdktypes.Coin, string, error) {

	dQ, err := QueryDelegation(address, valopAddress)
	if err != nil {
		return sdktypes.Coin{}, "", err
	}
	if dQ.Balance.Amount == "" || dQ.Delegation.Shares == "" {
		return sdktypes.Coin{}, "", fmt.Errorf("no delegation found")
	}
	balanceCoin, err := sdktypes.ParseCoinNormalized(dQ.Balance.Amount + dQ.Balance.Denom)
	if err != nil {
		return sdktypes.Coin{}, "", fmt.Errorf("cannot parse balance amount")
	}
	return balanceCoin, dQ.Delegation.Shares, nil
}
