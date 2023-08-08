package configs

import (
	"fmt"
	"load-balancer/commons"

	"github.com/spf13/viper"
)

var config commons.Config

func GetConfig() commons.Config {
	return config
}

func LoadConfig() (*commons.Config, error) {
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
