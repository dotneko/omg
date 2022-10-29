package main_test

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "runtime"
  "strings"
  "testing"
)

var (
  binName = "omg"
  fileName = ".omg.json"
)

func TestMain(m *testing.M) {
  fmt.Println("Building tool...")
  if runtime.GOOS == "windows" {
    binName += ".exe"
  }

  build := exec.Command("go", "build", "-o", binName)
  if err := build.Run(); err != nil {
    fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
    os.Exit(1)
  }
  fmt.Println("Running tests...")
  result := m.Run()
  fmt.Println("Cleaning up...")
  os.Remove(binName)
  os.Remove(fileName)
  os.Exit(result)
}

func TestOmgCLI(t *testing.T) {
  //a := omg.Account{Name: "Test1", Address: "Address1"}
  aNameAddress := "Test1 Address1"
  dir, err := os.Getwd()
  if err != nil {
    t.Fatal(err)
  }
  cmdPath := filepath.Join(dir, binName)

  t.Run("AddNewAccount", func(t *testing.T) {
    cmd := exec.Command(cmdPath, strings.Split("add " + aNameAddress, " ")...)
    if err := cmd.Run(); err != nil {
      t.Fatal(err)
    }
  })

  t.Run("ListWallets", func(t *testing.T) {
    cmd := exec.Command(cmdPath)
    out, err := cmd.CombinedOutput()
    if err != nil {
      t.Fatal(err)
    }
    
    expected := aNameAddress + "\n"

    if expected != string(out) {
      t.Errorf("Expected %q, got %q instead\n", expected, string(out))
    }
  })
}
