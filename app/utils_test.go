package app_test

import (
	"fmt"
	"testing"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"
	"github.com/shopspring/decimal"
)

func TestConversionDec(t *testing.T) {

	type conversionTest struct {
		inputAmt     decimal.Decimal
		inputDenom   string
		convertedAmt decimal.Decimal
	}
	d10 := decimal.NewFromFloat(10)
	d10e18 := decimal.NewFromFloat(10000000000000000000)

	var conversionTests = []conversionTest{
		{d10, cfg.Token, d10e18},
		{d10e18, cfg.BaseDenom, d10},
		{decimal.NewFromFloat(0.00000000012345), cfg.Token, decimal.NewFromFloat(123450000)},
		{decimal.NewFromFloat(123450000), cfg.BaseDenom, decimal.NewFromFloat(0.00000000012345)},
	}
	for _, test := range conversionTests {

		var convAmt decimal.Decimal
		if test.inputDenom == cfg.BaseDenom {
			convAmt = omg.DenomToTokenDec(test.inputAmt)
		} else if test.inputDenom == cfg.Token {
			convAmt = omg.TokenToDenomDec(test.inputAmt)
		}
		if !convAmt.Equal(test.convertedAmt) {
			t.Errorf("Expected %s, instead got %s", test.convertedAmt.String(), convAmt.String())
		}
	}
}

func TestStrSplitAmountDenomDec(t *testing.T) {
	type parseTest struct {
		inputStr    string
		parsedAmt   decimal.Decimal
		parsedDenom string
	}
	var tests = []parseTest{
		{"0nom", decimal.NewFromFloat(0), "nom"},
		{"1000nom", decimal.NewFromFloat(1000), "nom"},
		{"1000000anom", decimal.NewFromFloat(1000000), "anom"},
		{"-1234567890anom", decimal.NewFromFloat(-1234567890), "anom"},
		{fmt.Sprintf("0.00001%s", cfg.BaseDenom), decimal.NewFromFloat(0.00001), cfg.BaseDenom},
		{fmt.Sprintf("10%s", cfg.Token), decimal.NewFromFloat(10), cfg.Token},
	}
	for _, test := range tests {

		// Test parse string to amount/denom
		parsedAmt, parsedDenom, err := omg.StrSplitAmountDenomDec(test.inputStr)
		if err != nil {
			fmt.Println(err.Error())
			t.Errorf("Expected %s,%s; instead got %s,%s.\n", test.parsedAmt.String(), test.parsedDenom, parsedAmt.String(), parsedDenom)
		}
	}
}

func TestPrettifyDenom(t *testing.T) {
	type insertTest struct {
		amt decimal.Decimal
		out string
	}
	var testOutputs = []insertTest{
		{decimal.NewFromFloat(1), "1"},
		{decimal.NewFromFloat(999), "999"},
		{decimal.NewFromFloat(1000), "1_000"},
		{decimal.NewFromFloat(100099), "100_099"},
		{decimal.NewFromFloat(1000000), "1_000_000"},
		{decimal.NewFromFloat(99999999), "99_999_999"},
		{decimal.NewFromFloat(1000000000000000000), "1_000_000_000_000_000_000"},
	}
	for _, test := range testOutputs {
		output := omg.PrettifyDenom(test.amt)
		if test.out != output {
			t.Errorf("Expected %s; instead got %s\n", test.out, output)
		}
	}
}
