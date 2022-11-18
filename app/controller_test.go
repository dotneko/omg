package app_test

import (
	"fmt"
	"testing"

	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"
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