package app

import (
	"bufio"
	"fmt"
	"io"

	cfg "github.com/dotneko/omg/config"
)

func GetAliasAddress(r io.Reader, args ...string) (string, string, error) {
	var alias, address = "", ""
	if len(args) >= 2 {
		return args[0], args[1], nil
	}
	s := bufio.NewScanner(r)
	fmt.Print("Enter name/alias : ")
	s.Scan()
	if err := s.Err(); err != nil {
		return "", "", err
	}
	if len(s.Text()) == 0 {
		return "", "", fmt.Errorf("alias cannot be blank")
	}
	alias = s.Text()
	fmt.Print("Enter an address : ")
	s.Scan()
	if err := s.Err(); err != nil {
		return "", "", err
	}
	if len(s.Text()) == 0 {
		return "", "", fmt.Errorf("address cannot be blank")
	}
	address = s.Text()
	return alias, address, nil
}

// Get accounts for transaction from args or stdin
func GetTxAccounts(r io.Reader, action string, args ...string) (string, string, error) {
	var (
		acc1 string = ""
		acc2 string = ""
	)
	if len(args) >= 2 {
		acc1 = args[0]
		acc2 = args[1]
	}

	s := bufio.NewScanner(r)
	// Get input if no argument provided for 1st account
	if acc1 == "" {
		fmt.Printf("Enter account to %s from   : ", action)
		s.Scan()
		if err := s.Err(); err != nil {
			return "", "", err
		}
		if len(s.Text()) == 0 {
			return "", "", fmt.Errorf("alias cannot be blank")
		}
		acc1 = s.Text()
	}
	if acc2 == "" {
		if action == "delegate" {
			fmt.Print("Enter validator to delegate to : ")
		} else {
			fmt.Printf("Enter account to %s to : ", action)
		}
		s.Scan()
		if err := s.Err(); err != nil {
			return "", "", err
		}
		if len(s.Text()) == 0 {
			return "", "", fmt.Errorf("alias cannot be blank")
		}
		acc2 = s.Text()
	}
	return acc1, acc2, nil
}

// Get amount from stdin
func GetAmount(r io.Reader, action string, address string, args ...string) (float64, error) {
	var (
		amount      float64
		denom       string
		denomAmount float64
		err         error
	)
	if len(args) < 3 {
		// Prompt for amount if no amount provided in args
		s := bufio.NewScanner(r)
		fmt.Printf("Enter amount to %s : ", action)

		s.Scan()
		if err = s.Err(); err != nil {
			return 0, err
		}
		amount, denom, err = StrSplitAmountDenom(s.Text())
		if err != nil {
			return 0, err
		}
	} else {
		// Get amount from arguments
		amount, denom, err = StrSplitAmountDenom(args[2])
		if err != nil {
			return 0, err
		}
	}
	if amount == 0 {
		return 0, fmt.Errorf("amount cannot be 0")
	}
	// Convert to denom amount if token given
	if denom == cfg.Token {
		denomAmount = TokenToDenom(amount)
	} else if denom == cfg.Denom {
		denomAmount = amount
	} else {
		return 0, fmt.Errorf("invalid denomination - must be: %s / %s)", cfg.Denom, cfg.Token)
	}

	// Check if sufficient balance
	balance, err := GetBalanceAmount(address)
	if err != nil {
		return 0, err
	}
	if denomAmount > balance {
		return 0, fmt.Errorf("insufficient funds (requested %.0f%s)", denomAmount, cfg.Denom)
	}
	return denomAmount, nil
}
