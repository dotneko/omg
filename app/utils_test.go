package app_test

import (
	"strings"
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	omg "github.com/dotneko/omg/app"
	cfg "github.com/dotneko/omg/config"
)

func TestAmtToTokenDecCoin(t *testing.T) {
	type conversionTest struct {
		inputAmt     string
		convertedAmt sdktypes.DecCoin
	}
	d10 := "10.000000000000000000" + cfg.Token
	d10base := "10000000000000000000" + cfg.BaseDenom
	num1token := "0.000000000123450000" + cfg.Token
	num1base := "123450000" + cfg.BaseDenom

	var conversionTests = []conversionTest{
		{d10, sdktypes.NewDecCoin(cfg.Token, sdktypes.NewInt(10))},
		{d10base, sdktypes.NewDecCoin(cfg.Token, sdktypes.NewInt(10))},
		{num1token, sdktypes.NewDecCoinFromDec(cfg.Token, sdktypes.NewDecWithPrec(123450000, cfg.Decimals))},
		{num1base, sdktypes.NewDecCoinFromDec(cfg.Token, sdktypes.NewDecWithPrec(123450000, cfg.Decimals))},
	}
	for _, test := range conversionTests {
		var (
			convDecCoin sdktypes.DecCoin
			err         error
		)

		convDecCoin, err = omg.AmtToTokenDecCoin(test.inputAmt)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
		if !convDecCoin.IsEqual(test.convertedAmt) {
			t.Errorf("Expected %s, instead got %s", test.convertedAmt, convDecCoin)
		}
	}
}

func TestAmtToBaseCoin(t *testing.T) {
	type conversionTest struct {
		inputAmt     string
		convertedAmt sdktypes.Coin
	}
	d10 := "10.000000000000000000" + cfg.Token
	d10base := "10000000000000000000" + cfg.BaseDenom
	d10baseCoin, _ := sdktypes.ParseCoinNormalized(d10base)
	num1token := "0.000000000123450000" + cfg.Token
	num1base := "123450000" + cfg.BaseDenom
	num1baseCoin, _ := sdktypes.ParseCoinNormalized(num1base)

	var conversionTests = []conversionTest{
		{d10, d10baseCoin},
		{d10base, d10baseCoin},
		{num1token, num1baseCoin},
		{num1base, num1baseCoin},
	}
	for _, test := range conversionTests {
		var (
			convCoin sdktypes.Coin
			err      error
		)

		convCoin, err = omg.AmtToBaseCoin(test.inputAmt)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
		if !convCoin.IsEqual(test.convertedAmt) {
			t.Errorf("Expected %s, instead got %s", test.convertedAmt, convCoin)
		}
	}
}

func TestAmtToTokenStr(t *testing.T) {
	type conversionTest struct {
		inputAmt     string
		convertedAmt string
	}
	d10 := "10.000000000000000000" + cfg.Token
	d10base := "10000000000000000000" + cfg.BaseDenom
	num1token := "0.000000000123450000" + cfg.Token
	num1base := "123450000" + cfg.BaseDenom

	var conversionTests = []conversionTest{
		{d10, d10},
		{d10base, d10},
		{num1token, num1token},
		{num1base, num1token},
	}
	for _, test := range conversionTests {
		var (
			convAmt string
			err     error
		)

		convAmt = omg.AmtToTokenStr(test.inputAmt)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
		if !strings.EqualFold(test.convertedAmt, convAmt) {
			t.Errorf("Expected %s, instead got %s", test.convertedAmt, convAmt)
		}
	}
}

func TestAmtToBaseStr(t *testing.T) {
	type conversionTest struct {
		inputAmt     string
		convertedAmt string
	}
	d10 := "10.000000000000000000" + cfg.Token
	d10base := "10000000000000000000" + cfg.BaseDenom
	num1token := "0.000000000123450000" + cfg.Token
	num1base := "123450000" + cfg.BaseDenom

	var conversionTests = []conversionTest{
		{d10, d10base},
		{d10base, d10base},
		{num1token, num1base},
		{num1base, num1base},
	}
	for _, test := range conversionTests {
		var (
			convAmt string
			err     error
		)

		convAmt = omg.AmtToBaseStr(test.inputAmt)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
		if !strings.EqualFold(test.convertedAmt, convAmt) {
			t.Errorf("Expected %s, instead got %s", test.convertedAmt, convAmt)
		}
	}
}
func TestConvertAmt(t *testing.T) {

	type conversionTest struct {
		inputAmt     string
		inputDenom   string
		convertedAmt string
	}

	d10 := "10.000000000000000000" + cfg.Token
	d10base := "10000000000000000000" + cfg.BaseDenom
	num1token := "0.000000000123450000" + cfg.Token
	num1base := "123450000" + cfg.BaseDenom
	var conversionTests = []conversionTest{
		{d10, cfg.Token, d10},
		{d10, cfg.BaseDenom, d10base},
		{d10base, cfg.Token, d10},
		{d10base, cfg.BaseDenom, d10base},
		{num1token, cfg.Token, num1token},
		{num1token, cfg.BaseDenom, num1base},
		{num1token, cfg.Token, num1token},
		{num1base, cfg.BaseDenom, num1base},
	}
	for _, test := range conversionTests {

		convAmt := omg.ConvertAmt(test.inputAmt, test.inputDenom)

		if !strings.EqualFold(convAmt, test.convertedAmt) {
			t.Errorf("Convert %s to %s: Expected %s, instead got %s", test.inputAmt, test.inputDenom, test.convertedAmt, convAmt)
		}
	}
}

func TestStrSplitAmountDenom(t *testing.T) {
	type splitTest struct {
		input  string
		numstr string
		denom  string
	}

	d10 := "10.000000000000000000" + cfg.Token
	d10base := "10000000000000000000" + cfg.BaseDenom
	num1token := "0.000000000123450000" + cfg.Token
	num1base := "123450000" + cfg.BaseDenom
	var splitTests = []splitTest{
		{d10, "10.000000000000000000", cfg.Token},
		{d10base, "10000000000000000000", cfg.BaseDenom},
		{num1token, "0.000000000123450000", cfg.Token},
		{num1base, "123450000", cfg.BaseDenom},
	}
	for _, test := range splitTests {

		numstr, denom, err := omg.StrSplitAmountDenom(test.input)
		if err != nil {
			t.Errorf("Error testing %s: %s", test.input, err)
		}

		if !strings.EqualFold(numstr, test.numstr) {
			t.Errorf("Expected %s, instead got %s", test.numstr, numstr)
		}
		if !strings.EqualFold(denom, test.denom) {
			t.Errorf("Expected %s, instead got %s", test.denom, denom)
		}
	}
}

func TestPrettifyBaseAmt(t *testing.T) {
	type insertTest struct {
		amt string
		out string
	}
	var testOutputs = []insertTest{
		{"1" + cfg.BaseDenom, "1" + cfg.BaseDenom},
		{"999" + cfg.BaseDenom, "999" + cfg.BaseDenom},
		{"1000" + cfg.BaseDenom, "1_000" + cfg.BaseDenom},
		{"100099" + cfg.BaseDenom, "100_099" + cfg.BaseDenom},
		{"1000000" + cfg.BaseDenom, "1_000_000" + cfg.BaseDenom},
		{"99999999" + cfg.BaseDenom, "99_999_999" + cfg.BaseDenom},
		{"1000000000000000000" + cfg.BaseDenom, "1_000_000_000_000_000_000" + cfg.BaseDenom},
	}
	for _, test := range testOutputs {
		output := omg.PrettifyBaseAmt(test.amt)
		if test.out != output {
			t.Errorf("Expected %s; instead got %s\n", test.out, output)
		}
	}
}

func TestPrettifyTokenAmt(t *testing.T) {
	type insertTest struct {
		amt         string
		numDecimals int
		out         string
	}
	var testOutputs = []insertTest{
		{"0.1" + cfg.Token, -1, "0.1 " + cfg.Token},
		{"10.999" + cfg.Token, 4, "10.999 " + cfg.Token},
		{"10.999" + cfg.Token, 1, "10.9 " + cfg.Token},
		{"100.1000" + cfg.Token, 4, "100.1000 " + cfg.Token},
		{"1000.100099" + cfg.Token, 3, "1000.100 " + cfg.Token},
		{"1000000" + cfg.Token, 3, "1000000 " + cfg.Token},
		{"0.99999999" + cfg.Token, 6, "0.999999 " + cfg.Token},
		{"0.1000000000000000000" + cfg.Token, 20, "0.1000000000000000000 " + cfg.Token},
	}
	for _, test := range testOutputs {
		output := omg.PrettifyTokenAmt(test.amt, test.numDecimals)
		if test.out != output {
			t.Errorf("Expected %s; instead got %s\n", test.out, output)
		}
	}
}
