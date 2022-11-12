package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	config := ParseConfig()

	require.NotEmpty(t, config.App.Filename, "App:Filename should not be empty")
	require.NotEmpty(t, config.App.MinAliasLength, "App:MinAliasLength should not be empty")
	require.NotEmpty(t, config.Chain.Daemon, "Chain:Daemon should not be empty")
	require.NotEmpty(t, config.Chain.ChainId, "Chain:ChainId should not be empty")
	require.NotEmpty(t, config.Chain.AddressPrefix, "Chain:AddressPrefix should not be empty")
	require.NotEmpty(t, config.Chain.ValoperPrefix, "Chain:ValoperPrefix should not be empty")
	require.NotEmpty(t, config.Chain.AddressLength, "Chain:AddressLength should not be empty")
	require.NotEmpty(t, config.Chain.Denom, "Chain:Denom should not be empty")
	require.NotEmpty(t, config.Chain.Token, "Chain:Token should not be empty")
	require.NotEmpty(t, config.Chain.Decimals, "Chain:Decimals should not be empty")
	require.NotEmpty(t, config.Chain.DefaultFee, "Chain:DefaultFee should not be empty")
	require.NotEmpty(t, config.Chain.GasAdjust, "Chain:GasAdjust should not be empty")
	require.NotEmpty(t, config.Chain.KeyringDefault, "Chain:KeyringDefault should not be empty")
}
