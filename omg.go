package omg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Account struct {
	Name    string
	Address string
}

type Wallets []Account

// Checks if address is a validator account
func IsValidatorAddress(address string) bool {
	if address[:12] == "onomyvaloper" {
		return true
	}
	return false
}

// Checks if address is a wallet account
func IsNormalAddress(address string) bool {
	if address[:5] == "onomy" && address[5:12] != "valoper" {
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
		formatted += fmt.Sprintf("%2d: %10s [%s]\n", k, a.Name, a.Address)
	}
	return formatted
}

// Add method creates a new account and appends it to the list of Wallets
func (l *Wallets) Add(name string, address string) {
	a := Account{
		Name:    name,
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

func (l *Wallets) GetAddress(name string) string {
	for _, a := range *l {
		if name == a.Name {
			return a.Address
		}
	}
	return ""
}

// Save method encodes the Wallets list as JSON and saves it
func (l *Wallets) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, js, 0644)
}

// Get method opens the provided filename, decodes JSON and parses it to list
func (l *Wallets) Load(filename string) error {
	file, err := os.ReadFile(filename)
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
