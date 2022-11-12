package omg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Account struct {
	Alias   string
	Address string
}

type Wallets []Account

const (
	MinAliasLength int    = 3
	AddressPrefix  string = "onomy"
	ValoperPrefix  string = "onomyvaloper"
	bech32len      int    = 39
)

// Checks if address is a validator account
func IsValidatorAddress(address string) bool {
	if len(address) == (len(ValoperPrefix)+bech32len) && address[:12] == ValoperPrefix {
		return true
	}
	return false
}

// Checks if address is a wallet account
func IsNormalAddress(address string) bool {
	if len(address) == (len(AddressPrefix)+bech32len) && address[:5] == "onomy" {
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

// String prints out a formatted list
// Implements the fmt.Stringer interface
func (l *Wallets) String() string {
	formatted := ""
	for k, a := range *l {
		formatted += fmt.Sprintf("%2d: %10s [%s]\n", k, a.Alias, a.Address)
	}
	return formatted
}

// Add method creates a new account and appends it to the list of Wallets
func (l *Wallets) Add(alias string, address string) {
	a := Account{
		Alias:   alias,
		Address: address,
	}

	*l = append(*l, a)
}

// Delete method deletes an account from the list
func (l *Wallets) Delete(idx int) error {
	ls := *l
	if idx < 0 || idx >= len(ls) {
		return fmt.Errorf("Account at index %d does not exist", idx)
	}
	// Adjusting index for 0 based index
	*l = append(ls[:idx], ls[idx+1:]...)

	return nil
}

// Get address for wallet alias
func (l *Wallets) GetAddress(alias string) string {
	for _, a := range *l {
		if alias == a.Alias {
			return a.Address
		}
	}
	return ""
}

// Save method encodes the Wallets list as JSON and saves it
func (l *Wallets) Save(filealias string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return os.WriteFile(filealias, js, 0644)
}

// Get method opens the provided filealias, decodes JSON and parses it to list
func (l *Wallets) Load(filealias string) error {
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
