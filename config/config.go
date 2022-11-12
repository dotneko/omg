package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig   `mapstructure:"app"`
	Chain ChainConfig `mapstructure:"chain"`
}

type AppConfig struct {
	MinAliasLength int    `mapstructure:"alias_length"`
	Filename       string `mapstructure:"omg_filename"`
}

type ChainConfig struct {
	Daemon         string  `mapstructure:"daemon"`
	ChainId        string  `mapstructure:"chain_id"`
	AddressPrefix  string  `mapstructure:"address_prefix"`
	ValoperPrefix  string  `mapstructure:"valoper_prefix"`
	AddressLength  int     `mapstructure:"address_length"`
	Denom          string  `mapstructure:"denom"`
	Token          string  `mapstructure:"token"`
	Decimals       string  `mapstructure:"decimals"`
	DefaultFee     int     `mapstructure:"default_fee"`
	GasAdjust      float32 `mapstructure:"gas_adjust"`
	KeyringDefault string  `mapstructure:"keyring_backend"`
}

func ParseConfig() *Config {
	var config Config
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")

	home, _ := os.UserHomeDir()
	viper.AddConfigPath(home)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading configuration, %s", err.Error())
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
	return &config
}
