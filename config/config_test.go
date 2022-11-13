package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	err := ParseConfig("../")
	if err != nil {
		t.Error(err)
	}

	require.NotEmpty(t, OmgFilename, "OmgFilename should not be empty")
	require.NotEmpty(t, MinAliasLength, "MinAliasLength should not be empty")
	require.NotEmpty(t, Daemon, "Daemon should not be empty")
	require.NotEmpty(t, ChainId, "ChainId should not be empty")
	require.NotEmpty(t, AddressPrefix, "AddressPrefix should not be empty")
	require.NotEmpty(t, ValoperPrefix, "ValoperPrefix should not be empty")
	require.NotEmpty(t, Denom, "Denom should not be empty")
	require.NotEmpty(t, Token, "Token should not be empty")
	require.NotEmpty(t, Decimals, "Decimals should not be empty")
	require.NotEmpty(t, DefaultFee, "DefaultFee should not be empty")
	require.NotEmpty(t, GasAdjust, "GasAdjust should not be empty")
	require.NotEmpty(t, KeyringBackend, "KeyringDefault should not be empty")
}
