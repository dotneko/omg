package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"omg"
)

// Project defaults
var (
	daemon      = "onomyd"
	omgFileName = ".omg.json"
	chainId     = "onomy-testnet-1"
)

const (
	denom      = "anom"
	token      = "nom"
	decimals   = 18
	jsonFlag   = "-o json"
	defaultFee = 4998
	gasAdjust  = 1.5
	keyring    = "test"
)

// Types for JSON unmarshalling
type BalancesQuery struct {
	Balances   []DenomAmount
	Pagination PaginationStruct
}

type RewardsQuery struct {
	Rewards []ValidatorReward
	Total   []DenomAmount
}

type ValidatorReward struct {
	ValidatorAddress string `json:"validator_address"`
	Reward           []DenomAmount
}

type DenomAmount struct {
	Denom  string
	Amount string
}

type PaginationStruct struct {
	NextKey string `json:"next_key"`
	Total   string
}

func getNameAddress(r io.Reader, args ...string) (string, string, error) {
	var name, address = "", ""
	if len(args) >= 2 {
		return args[0], args[1], nil
	}
	s := bufio.NewScanner(r)
	fmt.Print("Enter acc name : ")
	s.Scan()
	if err := s.Err(); err != nil {
		return "", "", err
	}
	if len(s.Text()) == 0 {
		return "", "", fmt.Errorf("Name cannot be blank")
	}
	name = s.Text()
	fmt.Print("Enter address  : ")
	s.Scan()
	if err := s.Err(); err != nil {
		return "", "", err
	}
	if len(s.Text()) == 0 {
		return "", "", fmt.Errorf("Address cannot be blank")
	}
	address = s.Text()
	return name, address, nil
}

// Convert denom to token amount
func denomToToken(amt float64) float64 {
	return amt / math.Pow10(decimals)
}

func denomToStr(amt float64) string {
	return fmt.Sprintf("%.0f%s", amt, denom)
}
func tokenToDenom(amt float64) float64 {
	return amt * math.Pow10(decimals)
}

func strToFloat(amtstr string) (float64, error) {
	amt, err := strconv.ParseFloat(amtstr, 64)
	if err != nil {
		return -1, err
	}
	return amt, nil
}

// Parse balance
func getBalance(address string) (float64, error) {
	cmdStr := fmt.Sprintf("query bank balances %s %s", jsonFlag, address)
	out, err := exec.Command(daemon, strings.Split(cmdStr, " ")...).Output()
	if err != nil {
		return 0, err
	}
	if !json.Valid(out) {
		return 0, errors.New("Invalid json")
	}
	var b BalancesQuery
	if err = json.Unmarshal(out, &b); err != nil {
		return 0, err
	}
	balance, err := strconv.ParseFloat(b.Balances[0].Amount, 64)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

// Check balance method
func checkBalance(address string) {
	balance, _ := getBalance(address)
	fmt.Printf("Balance = %.0f anom (%8.5f nom)\n", balance, denomToToken(balance))
}

// Check reward method
func checkRewards(address string) {
	cmdStr := fmt.Sprintf("query distribution rewards %s %s", jsonFlag, address)

	out, err := exec.Command(daemon, strings.Split(cmdStr, " ")...).Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	if !json.Valid(out) {
		fmt.Println("Invalid json")
	}
	var r RewardsQuery
	if err = json.Unmarshal(out, &r); err != nil {
		fmt.Println(err)
	}
	for _, v := range r.Rewards {
		amt, err := strToFloat(v.Reward[0].Amount)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s - %8.5f nom\n", v.ValidatorAddress, denomToToken(amt))
	}
}

// Withdraw all rewards method
func withdrawRewards(name string, auto bool) {

	cmdStr := fmt.Sprintf("tx distribution withdraw-all-rewards --from %s", name)
	cmdStr += fmt.Sprintf(" --fees %d%s --gas auto --gas-adjustment %f", defaultFee, denom, gasAdjust)
	cmdStr += fmt.Sprintf(" --keyring-backend %s --chain-id %s", keyring, chainId)

	fmt.Printf("Executing: %s %s\n", daemon, cmdStr)
	cmd := exec.Command(daemon, strings.Split(cmdStr, " ")...)

	if auto {
		// Auto confirm transaction
		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		// Expect prompt and confirm with 'y'
		stdin.Write([]byte("y\n"))

		if err := cmd.Wait(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		// Interactive execution
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

}

// Get delegator and validator from args or stdin
func getDelegatorValidator(r io.Reader, args ...string) (string, string, error) {
	var (
		delegator string = ""
		validator string = ""
	)
	if len(args) >= 2 {
		delegator = args[0]
		validator = args[1]
	}

	s := bufio.NewScanner(r)
	// Get delegator input if no argument provided
	if delegator == "" {
		fmt.Print("Enter wallet to delegate from : ")
		s.Scan()
		if err := s.Err(); err != nil {
			return "", "", err
		}
		if len(s.Text()) == 0 {
			return "", "", fmt.Errorf("Name cannot be blank\n")
		}
		delegator = s.Text()
	}
	if validator == "" {
		fmt.Print("Enter validator to delegate to : ")
		s.Scan()
		if err := s.Err(); err != nil {
			return "", "", err
		}
		if len(s.Text()) == 0 {
			return "", "", fmt.Errorf("Validator cannot be blank\n")
		}
		validator = s.Text()
	}
	return delegator, validator, nil
}

// Get amount from stdin
func getDelegationAmount(r io.Reader, address string, args ...string) (float64, error) {
	var amount float64 = 0.0
	if len(flag.Args()) == 3 {
		amount, err := strconv.ParseFloat(flag.Args()[2], 64)
		if err != nil {
			fmt.Println(err)
		}
		return amount, nil
	}
	s := bufio.NewScanner(r)
	fmt.Printf("Enter amount to delegate [%s] : ", token)

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
		// Negative amounts represent approx remaining amount after delegation
		balance, err := getBalance(address)
		if err != nil {
			return 0, err
		}
		amount = balance + tokenToDenom(tokenAmt)
		if amount <= 0 {
			return 0, fmt.Errorf("Insufficient balance")
		}
		return amount, nil
	}
	amount = tokenToDenom(tokenAmt)
	return amount, nil
}

// Delegate to validator method
func delegateToValidator(delegator string, valAddress string, amount float64, auto bool) {
	// fmt.Printf("DelegateToValidator %s %s %s %t\n", delegator, valAddress, denomToStr(amount), auto)

	cmdStr := fmt.Sprintf("tx staking delegate %s %s --from %s", valAddress, denomToStr(amount), delegator)
	cmdStr += fmt.Sprintf(" --fees %d%s --gas auto --gas-adjustment %f", defaultFee, denom, gasAdjust)
	cmdStr += fmt.Sprintf(" --keyring-backend %s --chain-id %s", keyring, chainId)

	fmt.Printf("Executing: %s %s\n", daemon, cmdStr)
	cmd := exec.Command(daemon, strings.Split(cmdStr, " ")...)

	if auto {
		// Auto confirm transaction
		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		// Expect prompt and confirm with 'y'
		stdin.Write([]byte("y\n"))

		if err := cmd.Wait(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		// Interactive execution
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func main() {

	if os.Getenv("OMG_FILENAME") != "" {
		omgFileName = os.Getenv("OMG_FILENAME")
	}

	// Flag usage
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. ", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Usage infomation:\n")
		flag.PrintDefaults()
	}
	// Parsing command line flags
	add := flag.Bool("add", false, "Add [account_name] [address] to wallets list")
	auto := flag.Bool("auto", false, "Auto confirm transactions")
	balance := flag.String("balance", "", "Check bank balance for [account_name]")
	delegate := flag.Bool("delegate", false, "Delegate from [account_name] to [validator]")
	list := flag.Bool("list", false, "List all accounts")
	rewards := flag.String("rewards", "", "Check rewards for [account_name]")
	rm := flag.String("rm", "", "Remove [account_name] from list")
	wdall := flag.String("wdall", "", "Withdraw all rewards for [account_name]")

	flag.Parse()

	// Define an wallet list
	l := &omg.Wallets{}

	// Read items from file
	if err := l.Load(omgFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on number of arguments provided
	switch {

	case *add:
		name, address, err := getNameAddress(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Check if name exists
		existingAddress := ""
		for _, a := range *l {
			if a.Name == name {
				existingAddress = a.Address
			}
		}
		if existingAddress != "" {
			fmt.Printf("Aborting: %q already exists [%s]\n", name, existingAddress)
			os.Exit(1)
		}
		l.Add(name, address)
		// Save the new list
		if err := l.Save(omgFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Added %q, %q to wallets\n", name, address)

	case *balance != "":
		address := l.GetAddress(*balance)
		if address == "" {
			fmt.Printf("Error: account %q not found.\n", *balance)
			os.Exit(1)
		}
		checkBalance(address)

	case *delegate:
		delegator, validator, err := getDelegatorValidator(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Check if delegator in list and is not validator account
		delegatorAddress := l.GetAddress(delegator)
		if !omg.IsNormalAddress(delegatorAddress) {
			fmt.Errorf("Invalid delegator wallet: %s\n", delegatorAddress)
			os.Exit(1)
		}
		// Check if valid validator address
		valAddress := l.GetAddress(validator)
		if valAddress == "" {
			fmt.Errorf("Address not in list\n")
			os.Exit(1)
		}
		if !omg.IsValidatorAddress(valAddress) {
			fmt.Errorf("%q is not a validator address\n", valAddress)
			os.Exit(1)
		}
		// Check balance for delegator
		fmt.Printf("Delegator %s [%s]\n", delegator, delegatorAddress)
		checkBalance(delegatorAddress)
		
		amount, err := getDelegationAmount(os.Stdin, delegatorAddress, flag.Args()...)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		delegateToValidator(delegator, valAddress, amount, *auto)
	case *list:
		if len(*l) == 0 {
			fmt.Println("No accounts in store")
		} else {
			fmt.Print(l)
		}

	case *rewards != "":
		address := l.GetAddress(*rewards)
		if address == "" {
			fmt.Printf("Error: account %q not found.\n", *rewards)
			os.Exit(1)
		}
		checkRewards(address)

	case *rm != "":
		deleted := false
		for k, a := range *l {
			if *rm == a.Name {
				l.Delete(k)
				if err := l.Save(omgFileName); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				fmt.Printf("Deleted: %q [%q]\n", a.Name, a.Address)
				deleted = true
			}
		}
		if !deleted {
			fmt.Printf("%q not found.", *rm)
		}

	case *wdall != "":
		address := l.GetAddress(*wdall)
		if address == "" {
			fmt.Printf("Error: account %q not found.\n", *wdall)
			os.Exit(1)
		}
		withdrawRewards(*wdall, *auto)

	default:
		flag.Usage()
	}
}
