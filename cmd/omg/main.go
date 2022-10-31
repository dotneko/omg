package main

import (
	"flag"
	"fmt"
	"os"

	//"strings"

	"omg"
)

// Default filename
var omgFileName = ".omg.json"

func main() {

	if os.Getenv("ACC_FILENAME") != "" {
		omgFileName = os.Getenv("ACC_FILENAME")
	}

	// Flag usage
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. ", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Usage infomation:\n")
		flag.PrintDefaults()
	}
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
		if len(*l) == 0 {
			fmt.Println("No accounts in store")
		} else {
			fmt.Print(l)
		}

	case *add != "":
		// Get 1st string of remaining arguments as address
		if len(tail) == 0 {
			fmt.Println("Error: no address entered")
			os.Exit(0)
		}
		address := tail[0]
		l.Add(*add, address)
		// Save the new list
		if err := l.Save(omgFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Added %q, %q to wallets\n", *add, address)
	default:
		flag.Usage()
	}
}
