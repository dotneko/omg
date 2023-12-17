package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		OmgFilename    string `mapstructure:"addrbook_filename"`
		OmgPath        string `mapstructure:"addrbook_path"`
		MinAliasLength int    `mapstructure:"min_alias_length"`
	}
	Chain struct {
		Daemon        string `mapstructure:"daemon"`
		DaemonPath    string `mapstructure:"daemon_path"`
		ChainId       string `mapstructure:"chain_id"`
		AddressPrefix string `mapstructure:"address_prefix"`
		ValoperPrefix string `mapstructure:"valoper_prefix"`
		BaseDenom     string `mapstructure:"base_denom"`
		Token         string `mapstructure:"token"`
		Decimals      int64  `mapstructure:"decimals"`
	}
	Options struct {
		DefaultFee     string  `mapstructure:"default_fee"`
		GasAdjust      float32 `mapstructure:"gas_adjust"`
		KeyringBackend string  `mapstructure:"keyring_backend"`
		Remainder      string  `mapstructure:"remainder"`
	}
}

var (
	OmgFilepath    string
	MinAliasLength int
	Daemon         string
	ChainId        string
	AddressPrefix  string
	ValoperPrefix  string
	BaseDenom      string
	Token          string
	Decimals       int64
	DefaultFee     string
	GasAdjust      float32
	KeyringBackend string
	Remainder      string
)

func init() {
	err := ParseConfig("..")
	if err != nil {
		fmt.Println(err)
	}
}

func ParseConfig(pathstr string) error {
	var cfg Config

	viper.SetConfigName(".omgconfig.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

	// Check given pathstr if valid path
	_, err := os.Stat(pathstr)
	if err == nil {
		viper.AddConfigPath(pathstr)
	}
	// Check home directory
	home, _ := os.UserHomeDir()
	viper.AddConfigPath(home)

	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error reading configuration, %s", err.Error())
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	omgDir := cfg.App.OmgPath
	// Assign if envirnomental variable
	if len(omgDir) > 1 && omgDir[0] == '$' {
		omgDir = os.Getenv(omgDir[1:])
	}
	// Check if directory exists
	dir, err := os.Stat(omgDir)
	if err != nil {
		return fmt.Errorf("address book directory not found: %s", omgDir)
	}
	if !dir.IsDir() {
		return fmt.Errorf("configuration for address book - %q is not a directory", dir.Name())
	}

	OmgFilepath = filepath.Join(omgDir, cfg.App.OmgFilename)
	daemonDir := cfg.Chain.DaemonPath
	if len(daemonDir) > 1 && daemonDir[0] == '$' {
		daemonDir = os.Getenv(daemonDir[1:])
	}
	Daemon = filepath.Join(daemonDir, cfg.Chain.Daemon)
	// Check if file exists
	_, err = os.Stat(Daemon)
	if os.IsNotExist(err) {
		// Check for working daemon in path
		_, err = exec.Command("which", cfg.Chain.Daemon).Output()
		if err != nil {
			fmt.Printf("Error: cannot locate %s daemon.\n", cfg.Chain.Daemon)
			os.Exit(1)
		}
	}

	MinAliasLength = cfg.App.MinAliasLength
	if MinAliasLength == 0 {
		MinAliasLength = 3
	}
	ChainId = cfg.Chain.ChainId
	AddressPrefix = cfg.Chain.AddressPrefix
	ValoperPrefix = cfg.Chain.ValoperPrefix
	BaseDenom = cfg.Chain.BaseDenom
	Token = cfg.Chain.Token
	Decimals = cfg.Chain.Decimals
	DefaultFee = cfg.Options.DefaultFee
	GasAdjust = cfg.Options.GasAdjust
	KeyringBackend = cfg.Options.KeyringBackend
	Remainder = cfg.Options.Remainder

	return nil
}
