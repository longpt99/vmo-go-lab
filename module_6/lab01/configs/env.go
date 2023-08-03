package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Servers    []string `mapstructure:"servers"`
	RateLimit  int      `mapstructure:"rate_limit"`
	LogEnabled bool     `mapstructure:"log_enabled"`
	LogFile    string   `mapstructure:"log_file"`
	Database   struct {
		Redis struct {
			Db       int    `mapstructure:"db"`
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			Password string `mapstructure:"password"`
		} `mapstructure:"redis"`
	} `mapstructure:"database"`
}

var config Config

func GetConfig() Config {
	return config
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("system.yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		return nil, err
	}

	fmt.Println(config)
	return &config, nil
}
