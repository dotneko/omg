package app

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	cfg "github.com/dotneko/omg/config"
)

const (
	RAW    = "raw"
	DETAIL = "detail"
	HASH   = "hash"
	SHARES = "shares"
	TOKEN  = "token"
)

// Initialize denoms
func init() {
	tokenUnit := sdktypes.OneDec()
	baseUnit := sdktypes.NewDecWithPrec(1, cfg.Decimals) // 10^-cfgDecimals, e.g. 10^-18 (atto)

	err := sdktypes.RegisterDenom(cfg.Token, tokenUnit)
	if err != nil {
		fmt.Printf("Error registering %s: %s", cfg.Token, err)
		os.Exit(1)
	}

	err = sdktypes.RegisterDenom(cfg.BaseDenom, baseUnit)
	if err != nil {
		fmt.Printf("Error registering %s: %s", cfg.Token, err)
		os.Exit(1)
	}

}

// Convert base denom to token DecCoin amount
func AmtToTokenDecCoin(amt string) (sdktypes.DecCoin, error) {
	coin, err := sdktypes.ParseCoinNormalized(amt)
	if err != nil {
		return sdktypes.DecCoin{}, err
	}
	amtDecCoin, err := sdktypes.ConvertDecCoin(sdktypes.NewDecCoinFromCoin(coin), cfg.Token)
	return amtDecCoin, err
}

// Convert token to base Coin amount
// Wraps around ParseCoinNormalized
func AmtToBaseCoin(amt string) (sdktypes.Coin, error) {
	amtCoin, err := sdktypes.ParseCoinNormalized(amt)
	if err != nil {
		return sdktypes.Coin{}, err
	}
	return amtCoin, err
}

// Convert amount to token amount (string)
func AmtToTokenStr(amt string) string {
	tokenDecCoin, err := AmtToTokenDecCoin(amt)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	return tokenDecCoin.String()
}

// Convert amount to base denom anmount (string)
func AmtToBaseStr(amt string) string {
	amtCoin, err := sdktypes.ParseCoinNormalized(amt)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	return amtCoin.String()
}

// Old: Convert denom to token (Decimal); superceded by AmtToTokenDec
// func DenomToTokenDec(amt decimal.Decimal) decimal.Decimal {
// 	d := amt.Shift(-cfg.Decimals)
// 	return d
// }

// Convert Token to denom (Decimal)
// func TokenToDenomDec(amt decimal.Decimal) decimal.Decimal {
// 	d := amt.Shift(cfg.Decimals)
// 	return d
// }

// Convert amount to specified denom amount
func ConvertAmt(amt string, targetDenom string) string {
	if strings.EqualFold(targetDenom, cfg.Token) {
		// Convert to token decimal
		return AmtToTokenStr(amt)
	} else if strings.EqualFold(targetDenom, cfg.BaseDenom) {
		return AmtToBaseStr(amt)
	}
	return fmt.Sprintf("Error converting %s to %s", amt, targetDenom)
}

// Split denominated amount to amount and denom
func StrSplitAmountDenom(amtstr string) (string, string, error) {
	var NumericRegex = regexp.MustCompile(`[^0-9.-]+`)
	var AlphaRegex = regexp.MustCompile(`[^a-zA-z]+`)
	amtstr = strings.ReplaceAll(amtstr, "_", "")
	numstr := NumericRegex.ReplaceAllString(amtstr, "")
	denom := AlphaRegex.ReplaceAllString(amtstr, "")
	if strings.EqualFold(denom, cfg.BaseDenom) {
		return numstr, cfg.BaseDenom, nil
	} else if strings.EqualFold(denom, cfg.Token) {
		return numstr, cfg.Token, nil
	}
	return "", "", fmt.Errorf("denom must be %q or %q, got %q", cfg.BaseDenom, cfg.Token, denom)
}

func NormalizeAmountDenom(amtstr string) (string, error) {
	normalizedAmt, denom, err := StrSplitAmountDenom(amtstr)
	if err != nil {
		return "", err
	}
	return normalizedAmt + denom, nil
}

// Parse decimal string to fixed number of digits after the decimal point.
func PrettifyTokenAmt(amt string, numDecimals int) string {
	var (
		outStr string
	)
	numstr, denom, _ := StrSplitAmountDenom(amt)

	dotPos := strings.Index(numstr, ".")

	if dotPos == -1 || numDecimals == -1 {
		outStr = fmt.Sprintf("%s %s", numstr, denom)
	} else {
		endPos := dotPos + numDecimals + 1
		if endPos > len(numstr) {
			endPos = len(numstr)
		}
		outStr = fmt.Sprintf("%s %s", numstr[:endPos], denom)
	}
	return outStr
}

// Insert separator for non-decimal numbers as output
func PrettifyBaseAmt(amt string) string {
	var (
		amtStr string
		outStr string
	)
	amtCoin, err := sdktypes.ParseCoinNormalized(amt)
	if err != nil {
		return ""
	}
	_1000Coin := sdktypes.NewCoin(cfg.BaseDenom, sdktypes.NewInt(1000))
	if amtCoin.Amount.LT(_1000Coin.Amount) {
		return amt
	}

	amtStr = amtCoin.Amount.String()
	separator := "_"
	startIdx := len(amtStr) % 3
	if startIdx == 0 {
		startIdx = 3
	}
	outStr = amtStr[:startIdx]

	pos := startIdx
	for pos < len(amtStr) {
		outStr = outStr + separator + amtStr[pos:pos+3]
		pos = pos + 3
	}

	return outStr
}

func PrettifyAmount(amount string) string {
	amtCoin, err := sdktypes.ParseCoinNormalized(amount)
	if err != nil {
		return ""
	}
	if strings.EqualFold(amtCoin.Denom, cfg.BaseDenom) {
		return fmt.Sprintf("%s %s", PrettifyBaseAmt(amount), amtCoin.Denom)
	}
	if strings.EqualFold(amtCoin.Denom, cfg.Token) {
		return fmt.Sprintf("%s %s", amtCoin.Amount.String(), amtCoin.Denom)
	}
	return ""
}

func OutputAmount(out io.Writer, name, address string, baseAmount string, outType string) {
	amtCoin, err := sdktypes.ParseCoinNormalized(baseAmount)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !strings.EqualFold(amtCoin.Denom, cfg.BaseDenom) {
		fmt.Fprintf(out, "Warning: unexpected base denom: %q, expected %q", amtCoin.Denom, cfg.BaseDenom)
	}
	switch {
	case outType == RAW:
		fmt.Fprintf(out, "%s\n", baseAmount)
	case outType == TOKEN:
		tokenAmt, err := AmtToTokenDecCoin(baseAmount)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(out, "%s\n", tokenAmt)
	case outType == DETAIL:
		tokenAmt, err := AmtToTokenDecCoin(baseAmount)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(out, "%s [%s]\n> %s (%s)\n", name, address, tokenAmt.String(), amtCoin.String())
	default:
		tokenAmt, err := AmtToTokenDecCoin(baseAmount)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(out, "%20s : %20s (%s)\n", name, PrettifyTokenAmt(tokenAmt.String(), 4), PrettifyAmount(baseAmount))
	}
}
