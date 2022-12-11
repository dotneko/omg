package app

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	cfg "github.com/dotneko/omg/config"
	"github.com/shopspring/decimal"
)

const (
	RAW    = "raw"
	DETAIL = "detail"
	HASH   = "hash"
	SHARES = "shares"
	TOKEN  = "token"
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
	if strings.EqualFold(denom, cfg.BaseDenom) {
		convAmount = DenomToTokenDec(amount)
		convDenom = cfg.Token
	} else if strings.EqualFold(denom, cfg.Token) {
		convAmount = TokenToDenomDec(amount)
		convDenom = cfg.BaseDenom
	} else {
		fmt.Printf("Error: unrecognized denom - %q", denom)
		return decimal.NewFromInt(-1), ""
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
	if !strings.EqualFold(denom, cfg.BaseDenom) && !strings.EqualFold(denom, cfg.Token) {
		return decimal.NewFromInt(0), "", fmt.Errorf("denom must be %q or %q, got %q", cfg.BaseDenom, cfg.Token, denom)
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
	if strings.EqualFold(denom, cfg.BaseDenom) {
		return fmt.Sprintf("%s %s", PrettifyDenom(amount), denom)
	}
	if strings.EqualFold(denom, cfg.Token) {
		return fmt.Sprintf("%s %s", amount.String(), denom)
	}
	return ""
}

func OutputAmount(out io.Writer, name, address string, baseAmount decimal.Decimal, baseDenom, outType string) {
	if !strings.EqualFold(baseDenom, cfg.BaseDenom) {
		fmt.Fprintf(out, "Warning: unexpected base denom: %q, expected %q", baseDenom, cfg.BaseDenom)
	}
	switch {
	case outType == RAW:
		fmt.Fprintf(out, "%s%s\n", baseAmount.String(), baseDenom)
	case outType == TOKEN:
		fmt.Fprintf(out, "%s %s\n", DenomToTokenDec(baseAmount).StringFixed(18), cfg.Token)
	case outType == DETAIL:
		fmt.Fprintf(out, "%s [%s]\n> %s %s (%s%s)\n", name, address, DenomToTokenDec(baseAmount).String(), cfg.Token, baseAmount.String(), baseDenom)
	default:
		fmt.Fprintf(out, "%20s : %12s %s (%s%s)\n", name, DenomToTokenDec(baseAmount).StringFixed(4), cfg.Token, PrettifyDenom(baseAmount), cfg.BaseDenom)
	}
}
