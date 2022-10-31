package main

import (
	"flag"
	"fmt"
	"os"

	//"strings"

	"omg"
)

// Set filename
const omgFileName = ".omg.json"

func main() {
	// Parsing command line flags
	add := flag.String("add", "", "Add [account_name] [address] to wallets list")
	list := flag.Bool("list", false, "List all accounts")
	flag.Parse()
	tail := flag.Args()

	// Define an wallet list
	l := &omg.Wallets{}

	// Read items from file
	if err := l.Get(omgFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on number of arguments provided
	switch {
	case *list:
		count := 0
		for _, account := range *l {
			fmt.Println(account.Name, account.Address)
			count++
		}
		if count == 0 {
			fmt.Println("No accounts in store")
		}
	case *add != "":
		// Get 1st string of remaining arguments as address
		if len(tail) == 0 {
			fmt.Println("Error: no address entered")
			os.Exit(0)
		}
		address := tail[0]
		fmt.Printf("Address for %q : ", address)

		l.Add(*add, address)
		// Save the new list
		if err := l.Save(omgFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Added %q, %q to wallets\n", *add, address)
	default:
		fmt.Printf("Received %d args; Args: %q\n", len(os.Args), os.Args[1:])
	}
}
