package configs

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type File struct {
	Alias string `yaml:"alias"`
	Path  string `yaml:"path"`
}

type Service struct {
	Name           string `yaml:"name"`
	Host           string `yaml:"host"`
	User           string `yaml:"user"`
	Port           int    `yaml:"port"`
	PrivateKeyPath string `yaml:"private_key_path"`
	Krb5ConfPath   string `yaml:"krb5_conf_path"`
	Files          []File `yaml:"files"`
}

type Services struct {
	Services []Service `yaml:"services"`
}

type Config struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

func (s *Service) IsFullyConfigured() (bool, error) {
	propsDefined := s.Name != "" && s.Host != "" && s.User != "" && s.Port != 0 && (s.PrivateKeyPath != "" || s.Krb5ConfPath != "") && len(s.Files) > 0
	if !propsDefined {
		return false, errors.New("service properties are not fully defined")
	} else {
		if (s.PrivateKeyPath != "" && !fileExist(s.PrivateKeyPath)) ||
			(s.Krb5ConfPath != "" && !fileExist(s.Krb5ConfPath)) {
			return false, errors.New("private key or krb5.conf file does not exist")
		}

	}
	for _, file := range s.Files {
		if file.Path == "" || file.Alias == "" {
			return false, errors.New("file path or alias is not defined")

		}
	}
	return true, nil
}

func (s *Services) IsFullyConfigured() (bool, error) {
	for _, service := range s.Services {
		if _, err := service.IsFullyConfigured(); err != nil {
			return false, err
		}
	}

	return true, nil
}

func fileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func ParseFiles(files string) []File {
	var parsedFiles []File
	if files != "" {
		for index, file := range strings.Split(files, ",") {
			parsedFiles = append(parsedFiles, File{Alias: strconv.Itoa(index + 1), Path: file})
		}
	}
	return parsedFiles
}
