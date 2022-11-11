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
	"regexp"
	"strconv"
	"strings"
	"time"

	"omg"
)

// Project defaults
var (
	daemon      = "onomyd"
	omgFileName = ".omg.json"
	chainId     = "onomy-testnet-1"
)

const (
	denom          = "anom"
	token          = "nom"
	decimals       = 18
	jsonFlag       = "--output json"
	keyringFlag    = "--keyring-backend"
	defaultFee     = 4998
	gasAdjust      = 1.2
	keyringDefault = "test"
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
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type PaginationStruct struct {
	NextKey string `json:"next_key"`
	Total   string `json:"total"`
}

type KeysListQuery struct {
	Key []KeyStruct
}
type KeyStruct struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address"`
	Pubkey  string `json:"pubkey"`
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

// Convert denom to annotated string
func denomToStr(amt float64) string {
	return fmt.Sprintf("%.0f%s", amt, denom)
}

// Convert token amount to denom amount
func tokenToDenom(amt float64) float64 {
	return amt * math.Pow10(decimals)
}

// Strip non-numeric characters and convert to float
func strToFloat(amtstr string) (float64, error) {
	var nonNumericRegex = regexp.MustCompile(`[^0-9.]+`)
	numstr := nonNumericRegex.ReplaceAllString(amtstr, "")
	amt, err := strconv.ParseFloat(numstr, 64)
	if err != nil {
		return -1, err
	}
	return amt, nil
}

// Parse balance
func getBalances(address string) (float64, error) {
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

// Import address from keyring
func importFromKeyring(wallet *omg.Wallets, keyring string) (int, error) {
	cmdStr := fmt.Sprintf("keys list %s %s %s", jsonFlag, keyringFlag, keyring)
	out, err := exec.Command(daemon, strings.Split(cmdStr, " ")...).Output()
	if err != nil {
		return 0, err
	}
	if !json.Valid(out) {
		return 0, errors.New("Invalid json")
	}
	var k []KeyStruct
	if err = json.Unmarshal(out, &k); err != nil {
		return 0, err
	}
	if len(k) == 0 {
		fmt.Println("No addresses in keyring")
		return 0, nil
	}
	count := 0
	for _, key := range k {
		if wallet.GetAddress(key.Name) == "" {
			wallet.Add(key.Name, key.Address)
			count++
			fmt.Printf("Imported %s [%s]\n", key.Name, key.Address)
		} else {
			fmt.Printf("Skip existing key with name %q [%s]\n", key.Name, key.Address)
		}
	}
	return count, nil
}

// Check balance method
func checkBalances(address string) {
	balance, _ := getBalances(address)
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
func withdrawRewards(name string, keyring string, auto bool) {

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

// Get accounts for transaction from args or stdin
func getTxAccounts(r io.Reader, action string, args ...string) (string, string, error) {
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
			return "", "", fmt.Errorf("Name cannot be blank\n")
		}
		acc1 = s.Text()
	}
	if acc2 == "" {
		if action == "delegate" {
			fmt.Print("Enter validator to delegate to : ")
		} else {
			fmt.Print("Enter wallet to %s to : ", action)
		}
		s.Scan()
		if err := s.Err(); err != nil {
			return "", "", err
		}
		if len(s.Text()) == 0 {
			return "", "", fmt.Errorf("Name cannot be blank\n")
		}
		acc2 = s.Text()
	}
	return acc1, acc2, nil
}

// Get amount from stdin
func getAmount(r io.Reader, action string, address string, args ...string) (float64, error) {
	var amount float64 = 0.0
	balance, _ := getBalances(address)
	if len(flag.Args()) == 3 {
		// Check if denom appended
		argStr := flag.Args()[2]
		if len(argStr) > 5 && argStr[len(argStr)-4:] == denom {
			amount, err := strconv.ParseFloat(argStr[:len(argStr)-4], 64)
			if err != nil {
				return 0, err
			}
			fmt.Printf("Requested %.0f%s\n", amount, denom)
			if amount < 0 {
				// Negative amounts represent approx remaining amount after delegation
				if -amount > balance {
					return 0, fmt.Errorf("Error: insufficent funds (requested %.0f%s", amount+balance, denom)
				}
				return amount + balance, nil

			}
			if amount < 0 || amount > balance {
				return 0, fmt.Errorf("Error: insufficient funds (requested %.0f%s)", amount, denom)
			}
			return amount, nil
		}
		// If denom not included, treat as token amount
		amount, err := strconv.ParseFloat(flag.Args()[2], 64)
		if err != nil {
			fmt.Println(err)
		}
		amount = tokenToDenom(amount)
		if amount < 0 {
			// Negative amounts represent approx remaining amount after delegation
			if -amount > balance {
				return 0, fmt.Errorf("Error: insufficient funds (requested %.0f%s)", amount, denom)
			}
			return balance + amount, nil
		}
		if amount > balance {
			return 0, fmt.Errorf("Error: insufficient funds (requested %.0f%s)", amount, denom)
		}
		return amount, nil
	}
	s := bufio.NewScanner(r)
	fmt.Printf("Enter amount to %s [%s] : ", action, token)

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
		balance, err := getBalances(address)
		if err != nil {
			return 0, err
		}
		amount = balance + tokenToDenom(tokenAmt)
		if amount <= 0 {
			return 0, fmt.Errorf("Error: insufficient funds (requested %.0f%s)", amount, denom)
		}
		return amount, nil
	}
	amount = tokenToDenom(tokenAmt)
	return amount, nil
}

// Delegate to validator method
func delegateToValidator(delegator string, valAddress string, amount float64, keyring string, auto bool) {
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

// Send tokens between accounts method
func txSend(fromAddress string, toAddress string, amount float64, keyring string, auto bool) {
	// fmt.Printf("DelegateToValidator %s %s %s %t\n", delegator, valAddress, denomToStr(amount), auto)

	cmdStr := fmt.Sprintf("tx bank send %s %s %s", fromAddress, toAddress, denomToStr(amount))
	//cmdStr += fmt.Sprintf(" --fees %d%s --gas auto --gas-adjustment %f", defaultFee, denom, gasAdjust)
	cmdStr += fmt.Sprintf("--gas auto --gas-adjustment %f", gasAdjust)
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
			"%s: Onomy Manager. ", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Usage infomation:\n")
		flag.PrintDefaults()
	}
	// Parsing command line flags
	add := flag.Bool("add", false, "Add [account_name] [address] to wallets list")
	auto := flag.Bool("auto", false, "Auto confirm transaction flag")
	balances := flag.String("balances", "", "Check bank balances for [account_name]")
	convDenom := flag.String("convd", "", fmt.Sprintf("Convert (%s) to token (%s) amount", denom, token))
	convToken := flag.String("convt", "", fmt.Sprintf("Convert (%s) to denom (%s) amount", token, denom))
	delegate := flag.Bool("delegate", false, "Delegate from [account_name] to [validator_name]")
	importAddrs := flag.Bool("import", false, "Import addresses from keyring")
	keyring := flag.String("keyring", keyringDefault, "Keyring backend flag: e.g. test, pass")
	list := flag.Bool("list", false, "List all accounts")
	restake := flag.Bool("restake", false, "Restake from [account_name] to [validator]")
	rewards := flag.String("rewards", "", "Check rewards for [account_name]")
	rm := flag.String("rm", "", "Remove [account_name] from list")
	send := flag.Bool("send", false, "Send tokens from [from_account_name] to [to_account_name]")
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

	case *balances != "":
		address := l.GetAddress(*balances)
		if address == "" {
			fmt.Printf("Error: account %q not found.\n", *balances)
			os.Exit(1)
		}
		checkBalances(address)

	case *convDenom != "":
		amt, err := strToFloat(*convDenom)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("%f %s\n", denomToToken(amt), token)

	case *convToken != "":
		amt, err := strToFloat(*convToken)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("%.0f%s", tokenToDenom(amt), denom)

	case *delegate:
		delegator, validator, err := getTxAccounts(os.Stdin, "delegate", flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Check if delegator in list and is not validator account
		delegatorAddress := l.GetAddress(delegator)
		if delegatorAddress == "" {
			fmt.Println("Error: no delegator address")
		}
		if !omg.IsNormalAddress(delegatorAddress) {
			fmt.Fprintf(os.Stderr, "Invalid delegator wallet: %s\n", delegatorAddress)
			os.Exit(1)
		}
		// Check if valid validator address
		valAddress := l.GetAddress(validator)
		if valAddress == "" {
			fmt.Fprintf(os.Stderr, "Address not in list\n")
			os.Exit(1)
		}
		if !omg.IsValidatorAddress(valAddress) {
			fmt.Fprintf(os.Stderr, "%q is not a validator address\n", valAddress)
			os.Exit(1)
		}
		// Check balance for delegator
		fmt.Printf("Delegator %s [%s]\n", delegator, delegatorAddress)
		checkBalances(delegatorAddress)

		amount, err := getAmount(os.Stdin, "delegate", delegatorAddress, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		delegateToValidator(delegator, valAddress, amount, *keyring, *auto)

	case *importAddrs:
		num, err := importFromKeyring(l, *keyring)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if num > 0 {
			// Save the new list
			if err := l.Save(omgFileName); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
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

	case *restake:
		// Ensure all arguments provided
		if len(flag.Args()) != 2 {
			fmt.Printf("Insufficient arguments. Expecting [delgator] [valdiator]")
			os.Exit(1)
		}
		delegator := flag.Args()[0]
		validator := flag.Args()[1]
		delegatorAddress := l.GetAddress(delegator)
		if delegatorAddress == "" {
			fmt.Println("Error: no delegator address")
		}
		if !omg.IsNormalAddress(delegatorAddress) {
			fmt.Fprintf(os.Stderr, "Invalid delegator wallet: %s\n", delegatorAddress)
			os.Exit(1)
		}
		// Check if valid validator address
		valAddress := l.GetAddress(validator)
		if valAddress == "" {
			fmt.Fprintf(os.Stderr, "Address not in list\n")
			os.Exit(1)
		}
		if !omg.IsValidatorAddress(valAddress) {
			fmt.Fprintf(os.Stderr, "%q is not a validator address\n", valAddress)
			os.Exit(1)
		}
		// Check balance for delegator
		fmt.Printf("Delegator %s [%s]\n", delegator, delegatorAddress)
		balanceBefore, err := getBalances(delegatorAddress)
		if err != nil {
			fmt.Errorf("Error getting balance for %s\n", delegator)
			os.Exit(1)
		}
		fmt.Printf("Existing balance: %.0f %s\n", balanceBefore, denom)
		fmt.Printf("Withdrawing rewards for %s [%s]\n", delegator, delegatorAddress)
		withdrawRewards(delegator, *keyring, *auto)
		// Wait till balance is updated
		var balance *float64
		balance = new(float64)
		count := 0
		for count <= 10 {
			*balance, _ = getBalances(delegatorAddress)
			if *balance > balanceBefore {
				fmt.Printf("Updated balance  : %.0f %s\n", *balance, denom)
				break
			}
			count++
			time.Sleep(1 * time.Second)
		}
		// If balance not updated and -auto flag set then abort
		if *auto && *balance == balanceBefore {
			fmt.Printf("Error getting rewards. Aborting auto-restake")
			os.Exit(1)
		}
		// Restake amount leaving approx remainder of 1 token
		amount := *balance - tokenToDenom(1.0)
		delegateToValidator(delegator, valAddress, amount, *keyring, *auto)

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

	case *send:
		from, to, err := getTxAccounts(os.Stdin, "send", flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Check if delegator in list and is not validator account
		fromAddress := l.GetAddress(from)
		if fromAddress == "" {
			fmt.Println("Error: no from address")
		}
		if !omg.IsNormalAddress(fromAddress) {
			fmt.Fprintf(os.Stderr, "Invalid normal account: %s\n", fromAddress)
			os.Exit(1)
		}
		// Check if valid validator address
		toAddress := l.GetAddress(to)
		if toAddress == "" {
			fmt.Fprintf(os.Stderr, "Address not in list\n")
			os.Exit(1)
		}
		if !omg.IsNormalAddress(toAddress) {
			fmt.Fprintf(os.Stderr, "Invalid normal account: %s\n", toAddress)
			os.Exit(1)
		}
		// Check balance for delegator
		fmt.Printf("From: %s [%s]\n", from, fromAddress)
		checkBalances(fromAddress)

		amount, err := getAmount(os.Stdin, "send", fromAddress, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		txSend(from, to, amount, *keyring, *auto)

	case *wdall != "":
		address := l.GetAddress(*wdall)
		if address == "" {
			fmt.Printf("Error: account %q not found.\n", *wdall)
			os.Exit(1)
		}
		withdrawRewards(*wdall, *keyring, *auto)

	default:
		flag.Usage()
	}
}
