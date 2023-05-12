package configs

import (
	"github.com/a-dakani/LogSpy/logger"
	"gopkg.in/yaml.v3"
	"os"
)

func LoadServices(srvs *Services) {
	// Parse and validate services file
	yamlFile, err := os.ReadFile("config.services.yaml")
	if err != nil {
		logger.Fatal(err.Error())
	}
	if err := yaml.Unmarshal(yamlFile, &srvs); err != nil {
		logger.Fatal(err.Error())
	}
	if srvs.IsFullyConfigured() {
		logger.Info("config.services.yaml file loaded")
	} else {
		logger.Fatal("config.services.yaml file not fully configured")
	}
}

func LoadConfig(cfg *Config) {
	// Parse and validate config file
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		logger.Fatal(err.Error())
	}
	if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("config.yaml file loaded")

}
