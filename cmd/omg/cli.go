package cli

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"
)

// Types for JSON unmarshalling

func getAliasAddress(r io.Reader, args ...string) (string, string, error) {
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

// Import address from keyring
func importFromKeyring(l *omg.Accounts, keyring string) (int, error) {
	accounts, err := omg.GetKeyringAccounts(keyring)

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

// Check balance method
func checkBalances(address string) {
	balance, err := omg.GetBalanceAmount(address)
	if err != nil {
		fmt.Sprintln(err)
	}
	fmt.Printf("Balance = %.0f anom (%8.5f nom)\n", balance, omg.DenomToToken(balance))
}

// Check reward method
func checkRewards(address string) {
	r, err := omg.GetRewards(address)
	if err != nil {
		fmt.Sprintln(err)
	}
	for _, v := range r.Rewards {
		amt, err := omg.StrToFloat(v.Reward[0].Amount)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s - %8.5f nom\n", v.ValidatorAddress, omg.DenomToToken(amt))
	}
}

// Withdraw all rewards method
func withdrawRewards(alias string, keyring string, auto bool) {
	err := omg.TxWithdrawRewards(alias, keyring, auto)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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
func getAmount(r io.Reader, action string, address string, args ...string) (float64, error) {
	var amount float64 = 0.0
	balance, err := omg.GetBalanceAmount(address)

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
		amount = omg.TokenToDenom(amount)
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
		balance, err := omg.GetBalanceAmount(address)
		if err != nil {
			return 0, err
		}
		amount = balance + omg.TokenToDenom(tokenAmt)
		if amount <= 0 {
			return 0, fmt.Errorf("Error: insufficient funds (requested %.0f%s)", amount, cfg.Denom)
		}
		return amount, nil
	}
	amount = omg.TokenToDenom(tokenAmt)
	return amount, nil
}

// Delegate to validator
func delegateToValidator(delegator string, valAddress string, amount float64, keyring string, auto bool) {

	err := omg.TxDelegateToValidator(delegator, valAddress, amount, keyring, auto)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// Send tokens between accounts method
func txSend(fromAddress string, toAddress string, amount float64, keyring string, auto bool) {
	err := omg.TxSend(fromAddress, toAddress, amount, keyring, auto)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func Execute() {

	err := cfg.ParseConfig("../../")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannot load configuration file: %s", err.Error())
	}

	// Flag usage
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s - Onomy Manager. ", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Usage infomation:\n")
		flag.PrintDefaults()
	}
	// Parsing command line flags
	add := flag.Bool("add", false, "Add [alias] [address] to wallets list")
	auto := flag.Bool("auto", false, "Auto confirm transaction flag")
	balances := flag.String("balances", "", "Check bank balances for [alias]")
	convDenom := flag.String("convd", "", fmt.Sprintf("Convert (%s) to token (%s) amount", cfg.Denom, cfg.Token))
	convToken := flag.String("convt", "", fmt.Sprintf("Convert (%s) to denom (%s) amount", cfg.Token, cfg.Denom))
	delegate := flag.Bool("delegate", false, "Delegate from [alias] to [validator alias]")
	importAddrs := flag.Bool("import", false, "Import addresses from keyring")
	keyring := flag.String("keyring", cfg.KeyringBackend, "Keyring-backend flag: e.g. test, pass")
	list := flag.Bool("list", false, "List all accounts")
	restake := flag.Bool("restake", false, "Restake from [alias] to [validator alias]")
	rename := flag.Bool("rename", false, "Rename [alias] to [new alias]")
	rewards := flag.String("rewards", "", "Check rewards for [alias]")
	rm := flag.String("rm", "", "Remove [alias] from list")
	send := flag.Bool("send", false, "Send tokens from [alias from] to [alias to]")
	wdall := flag.String("wdall", "", "Withdraw all rewards for [alias]")

	flag.Parse()

	// Define an accounts list
	l := &omg.Accounts{}

	// Read from saved address book
	if err := l.Load(cfg.OmgFilename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {

	case *add:
		alias, address, err := getAliasAddress(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		err = l.Add(alias, address)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Save the new list
		if err := l.Save(cfg.OmgFilename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Added %q, %q to wallets\n", alias, address)

	case *balances != "":
		address := l.GetAddress(*balances)
		if address == "" {
			fmt.Printf("Error: account %q not found.\n", *balances)
			os.Exit(1)
		}
		checkBalances(address)

	case *convDenom != "":
		amt, err := omg.StrToFloat(*convDenom)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("%f %s\n", omg.DenomToToken(amt), cfg.Token)

	case *convToken != "":
		amt, err := omg.StrToFloat(*convToken)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("%.0f%s", omg.TokenToDenom(amt), cfg.Denom)

	case *delegate:
		delegator, validator, err := getTxAccounts(os.Stdin, "delegate", flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Check if delegator in list and is not validator account
		delegatorAddress := l.GetAddress(delegator)
		if delegatorAddress == "" {
			fmt.Fprintf(os.Stderr, "Error: no delegator address")
			os.Exit(1)
		}
		if !omg.IsNormalAddress(delegatorAddress) {
			fmt.Fprintf(os.Stderr, "Error: invalid delegator address: %s\n", delegatorAddress)
			os.Exit(1)
		}
		// Check if valid validator address
		valAddress := l.GetAddress(validator)
		if valAddress == "" {
			fmt.Fprintf(os.Stderr, "Error: no validator matching %q\n", validator)
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
			if err := l.Save(cfg.OmgFilename); err != nil {
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
			fmt.Printf("Error: insufficient arguments. Expecting [delgator] [valdiator]")
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
		balanceBefore, err := omg.GetBalanceAmount(delegatorAddress)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting balance for %s\n", delegator)
			os.Exit(1)
		}
		fmt.Printf("Existing balance: %.0f %s\n", balanceBefore, cfg.Denom)
		fmt.Printf("Withdrawing rewards for %s [%s]\n", delegator, delegatorAddress)
		withdrawRewards(delegator, *keyring, *auto)

		// Wait till balance is updated
		var balance *float64
		balance = new(float64)
		count := 0
		for count <= 10 {
			*balance, _ = omg.GetBalanceAmount(delegatorAddress)
			if *balance > balanceBefore {
				fmt.Printf("Updated balance  : %.0f %s\n", *balance, cfg.Denom)
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
		amount := *balance - omg.TokenToDenom(1.0)
		delegateToValidator(delegator, valAddress, amount, *keyring, *auto)

	case *rename:
		oldAlias := flag.Args()[0]
		newAlias := flag.Args()[1]

		if len(newAlias) < cfg.MinAliasLength {
			fmt.Fprintln(os.Stderr, "Error: Please use alias of at least 3 characters")
			os.Exit(1)
		}
		idx := l.GetIndex(oldAlias)
		err := l.Modify(idx, newAlias, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
			os.Exit(1)
		}
		err = l.Save(cfg.OmgFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Renamed: %q to %q\n", oldAlias, newAlias)

	case *rm != "":
		deleted := false
		for k, a := range *l {
			if *rm == a.Alias {
				l.DeleteIndex(k)
				if err := l.Save(cfg.OmgFilename); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				fmt.Printf("Deleted: %q [%q]\n", a.Alias, a.Address)
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
		fmt.Printf("To  : %s [%s]\n", to, toAddress)

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
