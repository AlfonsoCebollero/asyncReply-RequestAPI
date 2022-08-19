package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	AppConfig = &Config{}
)

type CadenceConfig struct {
	Domain    string            `yaml:"domain"`
	Service   string            `yaml:"service"`
	HostPort  string            `yaml:"hostPort"`
	Workflows map[string]string `yaml:"workflows"`
}

type Config struct {
	Env     string
	Cadence CadenceConfig
	Logger  *zap.Logger
}

// LoadConfiguration setup the config for the code run
func (h *Config) LoadConfiguration() {
	err := readConfigFile("appConfig", "yml")
	if err != nil {
		log.Panic("Could not read file appConfig")
	}
	if err = viper.Unmarshal(&h); err != nil {
		log.Panic(fmt.Sprintf("Unable to decode into struct, %v", err))
	}

	log.Info("Application configuration successfully loaded.")

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	h.Logger = logger

	logger.Debug("Finished loading Configuration!")
}
