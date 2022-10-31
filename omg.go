package omg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	//  "time"
)

type Account struct {
	Name    string
	Address string
}

type Wallets []Account

// String prints out a formatted list
// Implements the fmt.Stringer interface
func (l *Wallets) String() string {
	formatted := ""
	for k, a := range *l {
		formatted += fmt.Sprintf("%2d: %s [%s]\n", k, a.Name, a.Address)
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
	if idx <= 0 || idx > len(ls) {
		return fmt.Errorf("Account at index %d does not exist", idx)
	}
	// Adjusting index for 0 based index
	*l = append(ls[:idx-1], ls[idx:]...)

	return nil
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
func (l *Wallets) Get(filename string) error {
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
