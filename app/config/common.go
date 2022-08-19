package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var (
	configPath = "app/resources"
)

type Configuration interface {
	LoadConfiguration()
}

func readConfigFile(configName, configType string) error {
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.AutomaticEnv()
	viper.SetConfigType(configType)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		return err
	}

	return nil
}
