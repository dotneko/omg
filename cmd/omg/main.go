package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	//"strings"

	"omg"
)

// Default filename
var omgFileName = ".omg.json"

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
	del := flag.String("del", "", "Delete [account_name] from list")
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
			os.Exit(0)
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
	default:
		flag.Usage()
	}
}
