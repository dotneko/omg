package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		OmgPath        string `mapstructure:"addrbook_path"`
		OmgFilename    string `mapstructure:"addrbook_filename"`
		MinAliasLength int    `mapstructure:"alias_length"`
	}
	Chain struct {
		Daemon        string `mapstructure:"daemon"`
		ChainId       string `mapstructure:"chain_id"`
		AddressPrefix string `mapstructure:"address_prefix"`
		ValoperPrefix string `mapstructure:"valoper_prefix"`
		BaseDenom     string `mapstructure:"base_denom"`
		Token         string `mapstructure:"token"`
		Decimals      int32  `mapstructure:"decimals"`
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
	Decimals       int32
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
	if omgDir == "$HOME" {
		omgDir = home
	} else {
		dir, err := os.Stat(omgDir)
		if err != nil {
			return fmt.Errorf("path not found: %s", omgDir)
		}
		if !dir.IsDir() {
			return fmt.Errorf("%q is not a directory", dir.Name())
		}
	}
	OmgFilepath = filepath.Join(omgDir, cfg.App.OmgFilename)
	MinAliasLength = cfg.App.MinAliasLength
	Daemon = cfg.Chain.Daemon
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
