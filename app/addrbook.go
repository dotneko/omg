package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	cfg "github.com/dotneko/omg/config"
)

type Account struct {
	Alias   string
	Address string
}

type Accounts []Account

const (
	bech32len  int = 39
	AccNormal      = "normal"
	AccValoper     = "valoper"
)

// Checks if address is a validator account
func IsValidatorAddress(address string) bool {
	prefixLen := len(cfg.ValoperPrefix)
	if len(address) == (prefixLen+bech32len) && address[:prefixLen] == cfg.ValoperPrefix {
		return true
	}
	return false
}

// Checks if address is a wallet account
func IsNormalAddress(address string) bool {
	prefixLen := len(cfg.AddressPrefix)
	if len(address) == (prefixLen+bech32len) && address[:prefixLen] == cfg.AddressPrefix {
		return true
	}
	return false
}

// Checks if address is valid
func IsValidAddress(address string) bool {
	if !IsNormalAddress(address) && !IsValidatorAddress(address) {
		return false
	}
	return true
}

// Lists accounts
func (l *Accounts) String() string {
	formatted := ""
	for k, a := range *l {
		formatted += fmt.Sprintf("%2d: %15s [%s]\n", k, a.Alias, a.Address)
	}
	return formatted
}

func (l *Accounts) ListFiltered(accountType string, addressOnly bool) string {
	formatted := ""
	for k, a := range *l {
		include := true
		if accountType == AccNormal {
			include = IsNormalAddress(a.Address)
		} else if accountType == AccValoper {
			include = IsValidatorAddress(a.Address)
		}
		if include {
			if addressOnly {
				formatted += a.Address + "\n"
			} else {
				formatted += fmt.Sprintf("%2d: %15s [%s]\n", k, a.Alias, a.Address)
			}
		}
	}
	return formatted
}

// Add method creates a new account and appends it to the list of Accounts
func (l *Accounts) Add(alias string, address string) error {
	if !IsValidAddress(address) {
		return fmt.Errorf("%q is ot a valid address", address)
	}
	if existAddr := l.GetAddress(alias); existAddr != "" {
		return fmt.Errorf("existing entry for %s => %s", alias, existAddr)
	}
	a := Account{
		Alias:   alias,
		Address: address,
	}

	*l = append(*l, a)
	return nil
}

// Delete account from the list with matching index
func (l *Accounts) DeleteIndex(idx int) error {
	ls := *l
	if idx < 0 || idx >= len(ls) {
		return fmt.Errorf("Account at index %d does not exist", idx)
	}
	// Adjusting index for 0 based index
	*l = append(ls[:idx], ls[idx+1:]...)

	return nil
}

// Delete account from list with matching alias
func (l *Accounts) Delete(alias string) error {
	ls := *l
	idx := l.GetIndex(alias)

	if idx < 0 {
		return fmt.Errorf("Account matching %q does not exist", alias)
	}
	// Adjusting index for 0 based index
	*l = append(ls[:idx], ls[idx+1:]...)

	return nil
}

// Modify account details with matching index
func (l *Accounts) Modify(idx int, alias string, address string) error {
	ls := *l
	if idx < 0 || idx >= len(ls) {
		return fmt.Errorf("alias/address does not exist")
	}
	if alias != "" {
		ls[idx].Alias = alias
	}
	if IsValidAddress(address) {
		ls[idx].Address = address
	}
	*l = ls
	return nil
}

// Get address for given alias
func (l *Accounts) GetAddress(alias string) string {
	for _, a := range *l {
		if alias == a.Alias {
			return a.Address
		}
	}
	return ""
}

// Get index for given alias
func (l *Accounts) GetIndex(alias string) int {
	for idx, a := range *l {
		if alias == a.Alias {
			return idx
		}
	}
	return -1
}

// Save accounts in JSON format
func (l *Accounts) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, js, 0644)
}

// Load accounts from JSON file
func (l *Accounts) Load(filealias string) error {
	file, err := os.ReadFile(filealias)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if len(file) == 0 {
		return nil
	}
	return json.Unmarshal(file, l)
}

// Import addresses from keyring
func ImportFromKeyring(l *Accounts, keyring string) (int, error) {
	accounts, err := GetKeyringAccounts(keyring)

	if err != nil {
		return 0, err
	}

	count := 0
	for _, acc := range accounts {
		if l.GetAddress(acc.Alias) == "" {
			l.Add(acc.Alias, acc.Address)
			count++
			fmt.Printf("Imported %s [%s]\n", acc.Alias, acc.Address)
		} else {
			fmt.Printf("Skip existing key with alias %q [%s]\n", acc.Alias, acc.Address)
		}
	}
	return count, nil
}
