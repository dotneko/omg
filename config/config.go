package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	MinAliasLength int     `mapstructure:"alias_length"`
	OmgFilename    string  `mapstructure:"omg_filename"`
	Daemon         string  `mapstructure:"daemon"`
	ChainId        string  `mapstructure:"chain_id"`
	AddressPrefix  string  `mapstructure:"address_prefix"`
	ValoperPrefix  string  `mapstructure:"valoper_prefix"`
	BaseDenom      string  `mapstructure:"base_denom"`
	Token          string  `mapstructure:"token"`
	Decimals       int32   `mapstructure:"decimals"`
	DefaultFee     string  `mapstructure:"default_fee"`
	GasAdjust      float32 `mapstructure:"gas_adjust"`
	KeyringBackend string  `mapstructure:"keyring_backend"`
	Remainder      string  `mapstructure:"remainder"`
}

var (
	MinAliasLength int
	OmgFilename    string
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

	MinAliasLength = cfg.MinAliasLength
	OmgFilename = cfg.OmgFilename
	Daemon = cfg.Daemon
	ChainId = cfg.ChainId
	AddressPrefix = cfg.AddressPrefix
	ValoperPrefix = cfg.ValoperPrefix
	BaseDenom = cfg.BaseDenom
	Token = cfg.Token
	Decimals = cfg.Decimals
	DefaultFee = cfg.DefaultFee
	GasAdjust = cfg.GasAdjust
	KeyringBackend = cfg.KeyringBackend
	Remainder = cfg.Remainder

	return nil
}
