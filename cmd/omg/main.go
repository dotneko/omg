package main

import (
	"bufio"
	//"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
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
	jsonFlag   = "-o json"
	defaultFee = 4998
	gasAdjust  = 1.5
	keyring    = "test"
)

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

// Check balance method
func checkBalance(address string) {
	cmdStr := fmt.Sprintf("query bank balances %s %s", jsonFlag, address)
	cmd := exec.Command(daemon, strings.Split(cmdStr, " ")...)
	// if err := cmd.Run(); err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// }
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println(string(out))

}

// Check reward method
func checkRewards(address string) {
	cmdStr := fmt.Sprintf("query distribution rewards %s", address)
	cmd := exec.Command(daemon, strings.Split(cmdStr, " ")...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println(string(out))
}

// Withdraw all rewards method
func withdrawRewards(name string, auto bool) {

	cmdStr := fmt.Sprintf("tx distribution withdraw-all-rewards --from %s", name)
	cmdStr += fmt.Sprintf(" --fees %d%s --gas auto --gas-adjustment %f", defaultFee, denom, gasAdjust)
	cmdStr += fmt.Sprintf(" --keyring-backend %s --chain-id %s", keyring, chainId)

	fmt.Printf("Executing: %s %s\n", daemon, cmdStr)
	cmd := exec.Command(daemon, strings.Split(cmdStr, " ")...)

	if auto == true {
		// Auto confirm transaction
		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		// stdout, err := cmd.StdoutPipe()
		// if err != nil {
		// 	fmt.Fprintln(os.Stderr, err)
		// }
		// buf := bytes.NewBuffer(nil)
		// read stdout continuously in a separate go routine
		// go func() {
		// 	io.Copy(buf, stdout)
		// 	fmt.Fprint(os.Stdout, buf)
		// }()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		// Expect prompt to confirm with 'y'
		// Note need to fix missing output before the confirmation prompt
		stdin.Write([]byte("y\n"))
		//fmt.Fprint(os.Stdout, buf)
		// Parse output for error code
		// for _, line := range strings.Split(string(out), "\n") {
		// 	if strings.Contains(line, "Error") {
		// 		fmt.Fprintln(os.Stderr, line)
		// 	}
		// }
		if err := cmd.Wait(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		//fmt.Fprint(os.Stdout, buf)
	} else {
		// Interactive
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println("Done.")
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
	del := flag.String("del", "", "Delete [account_name] from list")
	bal := flag.String("bal", "", "Check bank balance for [account_name]")
	rewards := flag.String("rewards", "", "Check rewards for [account_name]")
	txwd := flag.String("txwd", "", "Withdraw all rewards for [account_name]")
	list := flag.Bool("list", false, "List all accounts")
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
	case *list:
		if len(*l) == 0 {
			fmt.Println("No accounts in store")
		} else {
			fmt.Print(l)
		}

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

	case *del != "":
		deleted := false
		for k, a := range *l {
			if *del == a.Name {
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
			fmt.Printf("%q not found.", *del)
		}
	case *rewards != "":
		address := l.GetAddress(*rewards)
		if address == "" {
			fmt.Printf("Error: account %q not found.\n", *rewards)
			os.Exit(1)
		}
		checkRewards(address)

	case *bal != "":
		address := l.GetAddress(*bal)
		if address == "" {
			fmt.Printf("Error: account %q not found.\n", *bal)
			os.Exit(1)
		}
		checkBalance(address)

	case *txwd != "":
		address := l.GetAddress(*txwd)
		if address == "" {
			fmt.Printf("Error: account %q not found.\n", *txwd)
			os.Exit(1)
		}
		withdrawRewards(*txwd, *auto)

	default:
		flag.Usage()
	}
}
