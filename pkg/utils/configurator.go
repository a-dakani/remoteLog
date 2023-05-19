package utils

import (
	"errors"
	"github.com/a-dakani/remoteLog/configs"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func LoadServices(srvs *configs.Services) error {

	// Just for Binary Usage get the config file path from executable path
	// Get the directory containing the program executable
	executablePath, err := os.Executable()
	if err != nil {
		return err
	}
	programDir := filepath.Dir(executablePath)

	// Parse and validate services file
	yamlFile, err := os.ReadFile(filepath.Join(programDir, "config.services.yaml"))
	//yamlFile, err := os.ReadFile(filepath.Join("config.services.yaml"))
	if err != nil {
		Fatal("config.services.yaml file not found")
		return err
	}
	if err = yaml.Unmarshal(yamlFile, &srvs); err != nil {
		Fatal("can not unmarshal config.services.yaml file")
		return err
	}

	if _, err = srvs.IsFullyConfigured(); err != nil {
		return errors.New("config.services.yaml file not fully configured")
	}

	Info("config.services.yaml file loaded")
	return nil
}

func LoadConfig(cfg *configs.Config) error {
	// Just for Binary Usage get the config file path from executable path
	// Get the directory containing the program executable
	executablePath, err := os.Executable()
	if err != nil {
		return err
	}
	programDir := filepath.Dir(executablePath)

	// Parse and validate config file
	yamlFile, err := os.ReadFile(filepath.Join(programDir, "config.yaml"))
	//yamlFile, err := os.ReadFile(filepath.Join("config.yaml"))
	if err != nil {
		Fatal("config.yaml file not found")
		return err
	}
	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		Fatal("can not unmarshal config.yaml file")
		return err
	}
	Info("config.yaml file loaded")
	return nil

}
