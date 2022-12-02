package app

import (
	"fmt"
	"regexp"
	"strings"

	cfg "github.com/dotneko/omg/config"
	"github.com/shopspring/decimal"
)

// Convert denom to token (Decimal)
func DenomToTokenDec(amt decimal.Decimal) decimal.Decimal {
	d := amt.Shift(-cfg.Decimals)
	return d
}

// Convert Token to denom (Decimal)
func TokenToDenomDec(amt decimal.Decimal) decimal.Decimal {
	d := amt.Shift(cfg.Decimals)
	return d
}

func ConvertDecDenom(amount decimal.Decimal, denom string) (decimal.Decimal, string) {
	var (
		convAmount decimal.Decimal
		convDenom  string
	)
	if denom == cfg.BaseDenom {
		convAmount = DenomToTokenDec(amount)
		convDenom = cfg.Token
	} else if denom == cfg.Token {
		convAmount = TokenToDenomDec(amount)
		convDenom = cfg.BaseDenom
	}
	return convAmount, convDenom
}

// Split denominated amount to amount and denom (Decimal)
func StrSplitAmountDenomDec(amtstr string) (decimal.Decimal, string, error) {
	var NumericRegex = regexp.MustCompile(`[^0-9.-]+`)
	var AlphaRegex = regexp.MustCompile(`[^a-zA-z]+`)
	amtstr = strings.ReplaceAll(amtstr, "_", "")
	numstr := NumericRegex.ReplaceAllString(amtstr, "")
	amt, err := decimal.NewFromString(numstr)
	if err != nil {
		return decimal.NewFromInt(0), "", err
	}
	denom := AlphaRegex.ReplaceAllString(amtstr, "")
	if denom != cfg.BaseDenom && denom != cfg.Token {
		return amt, "", nil
	}
	return amt, denom, nil
}

// Convert denom to annotated string
func DenomToStr(amt decimal.Decimal) string {
	return fmt.Sprintf("%s%s", amt.String(), cfg.BaseDenom)
}

// Strip non-numeric characters and convert to decimal
func StrToDec(amtstr string) (decimal.Decimal, error) {
	var NumericRegex = regexp.MustCompile(`[^0-9.]+`)
	numstr := NumericRegex.ReplaceAllString(amtstr, "")
	amt, err := decimal.NewFromString(numstr)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	return amt, nil
}

// Insert separator for non-decimal numbers as output
func PrettifyDenom(amt decimal.Decimal) string {
	var (
		amtStr string
		outStr string
		dotPos int
	)
	if amt.Abs().LessThan(decimal.NewFromInt(1000)) {
		return amt.String()
	}
	dotPos = strings.Index(amt.Abs().String(), ".")
	if dotPos == -1 {
		amtStr = amt.Abs().String()
	} else {
		s := strings.Split(amt.Abs().String(), ".")
		amtStr = s[0]
	}
	separator := "_"
	startIdx := len(amtStr) % 3
	if startIdx == 0 {
		startIdx = 3
	}
	outStr = amtStr[:startIdx]
	if amt.IsNegative() {
		outStr = "-" + outStr
	}
	pos := startIdx
	for pos < len(amtStr) {
		outStr = outStr + separator + amtStr[pos:pos+3]
		pos = pos + 3
	}
	if dotPos == -1 {
		return outStr
	}
	return outStr + "._"
}

func PrettifyAmount(amount decimal.Decimal, denom string) string {
	if denom == cfg.BaseDenom {
		return fmt.Sprintf("%s %s", PrettifyDenom(amount), denom)
	}
	if denom == cfg.Token {
		return fmt.Sprintf("%s %s", amount.String(), denom)
	}
	return ""
}