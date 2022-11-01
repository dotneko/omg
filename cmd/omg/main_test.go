package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName  = "omg"
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

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdPath := filepath.Join(dir, binName)
	aName1 := "Test1"
	aAddress1 := "TestAddress1"

	t.Run("AddNewAccountFromArguments", func(t *testing.T) {

		cmd := exec.Command(cmdPath, "-add", aName1, aAddress1)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	aName2 := "Test2"
	aAddress2 := "TestAddress2"

	t.Run("AddNewAccountFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(cmdStdin, aName2+"\n")
		io.WriteString(cmdStdin, aAddress2+"\n")
		cmdStdin.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("ListWallets", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("%2d: %s [%s]\n%2d: %s [%s]\n", 0, aName1, aAddress1, 1, aName2, aAddress2)

		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})
}
