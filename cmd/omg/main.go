package main

import (
  "fmt"
  "os"
  "strings"

  "omg"
)

// Set filename
const omgFileName = ".omg.json"

func main() {
  // Define an wallet list
  l := &omg.Wallets{}

  // Read items from file
  if err := l.Get(omgFileName); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }

  // Decide what to do based on number of arguments provided
  switch {
    case len(os.Args) == 2 && strings.ToLower(os.Args[1]) == "list":
      count := 0
      for _, account := range *l {
        fmt.Println(account.Name, account.Address)
        count++
      }
      if count == 0 {
        fmt.Println("No accounts in store")
      }
    case len(os.Args) == 4 && strings.ToLower(os.Args[1]) == "add":
        name := os.Args[2]
        address := os.Args[3]
        l.Add(name, address)
        // Save the new list
        if err := l.Save(omgFileName); err != nil {
          fmt.Fprintln(os.Stderr, err)
          os.Exit(1)
        }
        fmt.Printf("Added %q, %q to wallets\n", name, address)
    default:
      fmt.Printf("Received %d args; Args: %q\n", len(os.Args), os.Args[1:])
  }
}
