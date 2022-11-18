package app

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	cfg "github.com/dotneko/omg/config"
)

func GetAliasAddress(r io.Reader, args ...string) (string, string, error) {
	var alias, address = "", ""
	if len(args) >= 2 {
		return args[0], args[1], nil
	}
	s := bufio.NewScanner(r)
	fmt.Print("Enter acc alias : ")
	s.Scan()
	if err := s.Err(); err != nil {
		return "", "", err
	}
	if len(s.Text()) == 0 {
		return "", "", fmt.Errorf("Alias cannot be blank")
	}
	alias = s.Text()
	fmt.Print("Enter address  : ")
	s.Scan()
	if err := s.Err(); err != nil {
		return "", "", err
	}
	if len(s.Text()) == 0 {
		return "", "", fmt.Errorf("Address cannot be blank")
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
		fmt.Printf("Enter wallet to %s from : ", action)
		s.Scan()
		if err := s.Err(); err != nil {
			return "", "", err
		}
		if len(s.Text()) == 0 {
			return "", "", fmt.Errorf("Alias cannot be blank\n")
		}
		acc1 = s.Text()
	}
	if acc2 == "" {
		if action == "delegate" {
			fmt.Print("Enter validator to delegate to : ")
		} else {
			fmt.Printf("Enter wallet to %s to : ", action)
		}
		s.Scan()
		if err := s.Err(); err != nil {
			return "", "", err
		}
		if len(s.Text()) == 0 {
			return "", "", fmt.Errorf("Alias cannot be blank\n")
		}
		acc2 = s.Text()
	}
	return acc1, acc2, nil
}

// Get amount from stdin
func GetAmount(r io.Reader, action string, address string, args ...string) (float64, error) {
	var amount float64 = 0.0
	balance, err := GetBalanceAmount(address)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if len(flag.Args()) == 3 {
		// Check if denom appended
		argStr := flag.Args()[2]
		if len(argStr) > 5 && argStr[len(argStr)-4:] == cfg.Denom {
			amount, err := strconv.ParseFloat(argStr[:len(argStr)-4], 64)
			if err != nil {
				return 0, err
			}
			fmt.Printf("Requested %.0f%s\n", amount, cfg.Denom)
			if amount < 0 {
				// Negative amounts represent approx remaining amount after delegation
				if -amount > balance {
					return 0, fmt.Errorf("Error: insufficent funds (requested %.0f%s", amount+balance, cfg.Denom)
				}
				return amount + balance, nil

			}
			if amount < 0 || amount > balance {
				return 0, fmt.Errorf("Error: insufficient funds (requested %.0f%s)", amount, cfg.Denom)
			}
			return amount, nil
		}
		// If denom not included, treat as token amount
		amount, err := strconv.ParseFloat(flag.Args()[2], 64)
		if err != nil {
			fmt.Println(err)
		}
		amount = TokenToDenom(amount)
		if amount < 0 {
			// Negative amounts represent approx remaining amount after delegation
			if -amount > balance {
				return 0, fmt.Errorf("Error: insufficient funds (requested %.0f%s)", amount, cfg.Denom)
			}
			return balance + amount, nil
		}
		if amount > balance {
			return 0, fmt.Errorf("Error: insufficient funds (requested %.0f%s)", amount, cfg.Denom)
		}
		return amount, nil
	}
	s := bufio.NewScanner(r)
	fmt.Printf("Enter amount to %s [%s] : ", action, cfg.Token)

	s.Scan()
	if err := s.Err(); err != nil {
		return 0, err
	}
	if len(s.Text()) == 0 {
		return 0.0, fmt.Errorf("Invalid amount")
	}

	tokenAmt, err := strconv.ParseFloat(s.Text(), 64)
	if err != nil {
		return 0, err
	}
	if tokenAmt == 0 {
		return 0, fmt.Errorf("Invalid amount %f", amount)
	}
	if tokenAmt < 0 {
		// Negative amounts represent approx remaining amount after action
		balance, err := GetBalanceAmount(address)
		if err != nil {
			return 0, err
		}
		amount = balance + TokenToDenom(tokenAmt)
		if amount <= 0 {
			return 0, fmt.Errorf("Error: insufficient funds (requested %.0f%s)", amount, cfg.Denom)
		}
		return amount, nil
	}
	amount = TokenToDenom(tokenAmt)
	return amount, nil
}
