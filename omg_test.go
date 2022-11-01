package omg_test

import (
	"os"
	"testing"

	"omg"
)

func TestAdd(t *testing.T) {
	l := omg.Wallets{}

	name := "Test1"
	address := "onomy123456890111111111111111111111111111111"
	l.Add(name, address)

	if l[0].Name != name {
		t.Errorf("Expected %q, got %q instead.", name, l[0].Name)
	}
	if l[0].Address != address {
		t.Errorf("Expected %q, got %q instead.", address, l[0].Address)
	}
}

func TestDelete(t *testing.T) {
	l := omg.Wallets{}

	accounts := []omg.Account{
		{Name: "Test1", Address: "onomy123456890111111111111111111111111111111"},
		{Name: "Test2", Address: "onomy123456890222222222222222222222222222222"},
		{Name: "Test3", Address: "onomy123456890333333333333333333333333333333"},
	}

	for _, a := range accounts {
		l.Add(a.Name, a.Address)
	}
	if l[0].Name != accounts[0].Name {
		t.Errorf("Expected %q, got %q instead.", accounts[0].Name, l[0].Name)
	}
	l.Delete(2)

	if len(l) != 2 {
		t.Errorf("Expected list length %d, got %d instead.", 2, len(l))
	}

	if l[1] != accounts[1] {
		t.Errorf("Expected %q, got %q instead.", accounts[1], l[1])
	}
}

func TestSaveLoad(t *testing.T) {
	l1 := omg.Wallets{}
	l2 := omg.Wallets{}

	account := omg.Account{Name: "Test1", Address: "onomy123456890111111111111111111111111111111"}

	l1.Add(account.Name, account.Address)
	if l1[0] != account {
		t.Errorf("Expected %q, got %q instead.", account, l1[0])
	}

	tf, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
	defer os.Remove(tf.Name())

	if err := l1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving list to file: %s", err)
	}
	if err := l2.Load(tf.Name()); err != nil {
		t.Fatalf("Error getting list from file: %s", err)
	}
	if l1[0].Name != l2[0].Name {
		t.Errorf("Name %q should match name %q", l1[0].Name, l2[0].Name)
	}
}
