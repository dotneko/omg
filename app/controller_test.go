package app_test

import (
	"fmt"
	"testing"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"
	"github.com/shopspring/decimal"
)

func TestConversion(t *testing.T) {

	type conversionTest struct {
		inputAmt     float64
		inputDenom   string
		convertedAmt float64
	}
	var conversionTests = []conversionTest{

		{10, cfg.Token, 10000000000000000000},
		{10000000000000000000, cfg.Denom, 10},
	}
	for _, test := range conversionTests {

		var convAmt float64
		if test.inputDenom == cfg.Denom {
			convAmt = omg.DenomToToken(test.inputAmt)
		} else if test.inputDenom == cfg.Token {
			convAmt = omg.TokenToDenom(test.inputAmt)
		}
		if convAmt != test.convertedAmt {
			t.Errorf("Expected %f, instead got %f", test.convertedAmt, convAmt)
		}
	}
}

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
		{d10e18, cfg.Denom, d10},
		{decimal.NewFromFloat(0.00000000012345), cfg.Token, decimal.NewFromFloat(123450000)},
		{decimal.NewFromFloat(123450000), cfg.Denom, decimal.NewFromFloat(0.00000000012345)},
	}
	for _, test := range conversionTests {

		var convAmt decimal.Decimal
		if test.inputDenom == cfg.Denom {
			convAmt = omg.DenomToTokenDec(test.inputAmt)
		} else if test.inputDenom == cfg.Token {
			convAmt = omg.TokenToDenomDec(test.inputAmt)
		}
		if !convAmt.Equal(test.convertedAmt) {
			t.Errorf("Expected %s, instead got %s", test.convertedAmt.String(), convAmt.String())
		}
	}
}

func TestStrSplitAmountDenom(t *testing.T) {
	type parseTest struct {
		inputStr    string
		parsedAmt   float64
		parsedDenom string
	}
	var tests = []parseTest{
		{"0nom", 0, "nom"},
		{"1000nom", 1000, "nom"},
		{"1000000anom", 1000000, "anom"},
		{"-1234567890anom", -1234567890, "anom"},
		{fmt.Sprintf("0.00001%s", cfg.Denom), 0.00001, cfg.Denom},
		{fmt.Sprintf("10%s", cfg.Token), 10, cfg.Token},
		{"0.123", 0.123, ""},
		{"99999", 99999, ""},
		{"-100", -100, ""},
		{"-1", -1, ""},
	}
	for _, test := range tests {

		// Test parse string to amount/denom
		parsedAmt, parsedDenom, err := omg.StrSplitAmountDenom(test.inputStr)
		if err != nil {
			fmt.Println(err.Error())
			t.Errorf("Expected %f,%s; instead got %f,%s.\n", test.parsedAmt, test.parsedDenom, parsedAmt, parsedDenom)
		}
	}
}

func TestPrettifyDenom(t *testing.T) {
	type insertTest struct {
		amt float64
		out string
	}
	var testOutputs = []insertTest{
		{1, "1"},
		{999, "999"},
		{1000, "1,000"},
		{100099, "100,099"},
		{1000000, "1,000,000"},
		{99999999, "99,999,999"},
		{1000000000000000000, "1,000,000,000,000,000,000"},
	}
	for _, test := range testOutputs {
		output := omg.PrettifyDenom(test.amt)
		if test.out != output {
			t.Errorf("Expected %s; instead got %s\n", test.out, output)
		}
	}
}
