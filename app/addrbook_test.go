package app_test

import (
	"fmt"
	"os"
	"testing"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"
)

func init() {
	err := cfg.ParseConfig("../")
	if err != nil {
		fmt.Println(err.Error())
	}
}
func TestAddressTypes(t *testing.T) {
	normalAddress := "onomy123456890111111111111111111111111111111"
	validatorAddress := "onomyvaloper123456890111111111111111111111111111111"
	invalidPrefix := "cosmosnotavalidaddress0000000000000000000000"
	invalidLength := "onomy1234568901111111111111111111111111111119"
	type typeTest struct {
		address     string
		isNormal    bool
		isValidator bool
		isValid     bool
	}
	var typeTests = []typeTest{
		{normalAddress, true, false, true},
		{validatorAddress, false, true, true},
		{invalidPrefix, false, false, false},
		{invalidLength, false, false, false},
	}
	for _, test := range typeTests {
		if omg.IsNormalAddress(test.address) != test.isNormal {
			t.Errorf("Expected %q isNormal = %t", test.address, test.isNormal)
		}
		if omg.IsValidatorAddress(test.address) != test.isValidator {
			t.Errorf("Expected %q isValidator = %t", test.address, test.isValidator)
		}
		if omg.IsValidAddress(test.address) != test.isValid {
			t.Errorf("Expected %q isValid = %t", test.address, test.isValid)
		}
	}
}

func TestAdd(t *testing.T) {
	l := omg.Accounts{}

	alias := "Test1"
	address := "onomy123456890111111111111111111111111111111"
	l.Add(alias, address)

	if l[0].Alias != alias {
		t.Errorf("Expected %q, got %q instead.", alias, l[0].Alias)
	}
	if l[0].Address != address {
		t.Errorf("Expected %q, got %q instead.", address, l[0].Address)
	}
}

func TestDeleteIndex(t *testing.T) {
	l := omg.Accounts{}

	accounts := []omg.Account{
		{Alias: "Test1", Address: "onomy123456890111111111111111111111111111111"},
		{Alias: "Test2", Address: "onomy123456890222222222222222222222222222222"},
		{Alias: "Test3", Address: "onomy123456890333333333333333333333333333333"},
	}

	for _, a := range accounts {
		l.Add(a.Alias, a.Address)
	}
	if l[0].Alias != accounts[0].Alias {
		t.Errorf("Expected %q, got %q instead.", accounts[0].Alias, l[0].Alias)
	}

	l.DeleteIndex(2)
	if len(l) != 2 {
		t.Errorf("Expected list length %d, got %d instead.", 2, len(l))
	}

	if l[1] != accounts[1] {
		t.Errorf("Expected %q, got %q instead.", accounts[1], l[1])
	}
}
func TestDelete(t *testing.T) {
	l := omg.Accounts{}

	accounts := []omg.Account{
		{Alias: "Test1", Address: "onomy123456890111111111111111111111111111111"},
		{Alias: "Test2", Address: "onomy123456890222222222222222222222222222222"},
		{Alias: "Test3", Address: "onomy123456890333333333333333333333333333333"},
	}

	for _, a := range accounts {
		l.Add(a.Alias, a.Address)
	}
	if l[0].Alias != accounts[0].Alias {
		t.Errorf("Expected %q, got %q instead.", accounts[0].Alias, l[0].Alias)
	}

	l.Delete("Test1")
	if len(l) != 2 {
		t.Errorf("Expected list length %d, got %d instead.", 2, len(l))
	}

	if l[0] != accounts[1] {
		t.Errorf("Expected %q, got %q instead.", accounts[1], l[0])
	}

}

func TestModify(t *testing.T) {
	l := omg.Accounts{}

	accounts := []omg.Account{
		{Alias: "Test1", Address: "onomy123456890111111111111111111111111111111"},
		{Alias: "Test2", Address: "onomy123456890222222222222222222222222222222"},
		{Alias: "Test3", Address: "onomy123456890333333333333333333333333333333"},
	}

	for _, a := range accounts {
		l.Add(a.Alias, a.Address)
	}
	if l[0].Alias != accounts[0].Alias {
		t.Errorf("Expected %q, got %q instead.", accounts[0].Alias, l[0].Alias)
	}

	newAlias := "Modified2"
	newAddress := "onomy098765432111111111111111111111111111992"
	invalidAddress := "onomy098765432111111111111111111111111111992invalid"
	l.Modify(1, newAlias, newAddress)

	if l[1].Alias != newAlias {
		t.Errorf("Expected %q, got %q instead.", accounts[1].Alias, newAlias)
	}
	if l[1].Address != newAddress {
		t.Errorf("Expected %q, got %q instead.", accounts[1].Address, newAddress)
	}
	l.Modify(1, newAlias, invalidAddress)
	if l[1].Address == invalidAddress {
		t.Errorf("Expected %q, got %q instead.", accounts[1].Address, invalidAddress)
	}

}
func TestSaveLoad(t *testing.T) {
	l1 := omg.Accounts{}
	l2 := omg.Accounts{}

	account := omg.Account{Alias: "Test1", Address: "onomy123456890111111111111111111111111111111"}

	l1.Add(account.Alias, account.Address)
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
	if l1[0].Alias != l2[0].Alias {
		t.Errorf("Alias %q should match alias %q", l1[0].Alias, l2[0].Alias)
	}
}
