package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	mockBinDir = "path/to/bin"
	mockBin    = "onomyd"
)

func SetupMockFiles() {
	// Setup mock binary for testing
	// Assumes config specifies binary in path ./path/to/bin
	if err := os.MkdirAll(mockBinDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	binfile, err := os.Create(filepath.Join(".", mockBinDir, mockBin))
	if err != nil {
		log.Fatal(err)
	}
	binfile.Close()
}

func CleanupMockFiles() {
	rootDir := strings.Split(mockBinDir, "/")[0]
	if err := os.RemoveAll(rootDir); err != nil {
		log.Fatal(err)
	}
}
func TestParseConfig(t *testing.T) {

	SetupMockFiles()
	defer CleanupMockFiles()

	err := ParseConfig("..")
	if err != nil {
		t.Error(err)
	}

	require.NotEmpty(t, OmgFilepath, "OmgFilename should not be empty")
	fmt.Printf("OmgFilepath = %s\n", OmgFilepath)
	require.NotEmpty(t, MinAliasLength, "MinAliasLength should not be empty")
	require.NotEmpty(t, Daemon, "Daemon should not be empty")
	fmt.Printf("DaemonFilepath = %s\n", Daemon)
	require.NotEmpty(t, ChainId, "ChainId should not be empty")
	require.NotEmpty(t, AddressPrefix, "AddressPrefix should not be empty")
	require.NotEmpty(t, ValoperPrefix, "ValoperPrefix should not be empty")
	require.NotEmpty(t, BaseDenom, "BaseDenom should not be empty")
	require.NotEmpty(t, Token, "Token should not be empty")
	require.NotEmpty(t, Decimals, "Decimals should not be empty")
	require.NotEmpty(t, DefaultFee, "DefaultFee should not be empty")
	require.NotEmpty(t, GasAdjust, "GasAdjust should not be empty")
	require.NotEmpty(t, KeyringBackend, "KeyringDefault should not be empty")
	require.NotEmpty(t, Remainder, "Remainder should not be empty")

}
