package configs

type Service struct {
	Name           string   `yaml:"name"`
	Host           string   `yaml:"host"`
	User           string   `yaml:"user"`
	Port           int      `yaml:"port"`
	PrivateKeyPath string   `yaml:"private_key_path"`
	Files          []string `yaml:"files"`
}

type Services struct {
	Services []Service `yaml:"services"`
}

type Config struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

func (s *Service) IsFullyConfigured() bool {
	// TODO: check if private key exists
	// TODO: check if files exist
	return s.Name != "" && s.Host != "" && s.User != "" && s.Port != 0 && s.PrivateKeyPath != "" && len(s.Files) > 0
}
func (s *Services) IsFullyConfigured() bool {
	for _, service := range s.Services {
		if !service.IsFullyConfigured() {
			return false
		}
	}
	return true
}
