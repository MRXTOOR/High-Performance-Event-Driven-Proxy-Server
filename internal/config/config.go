package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ListenConfig struct {
	HTTP string `yaml:"http"`
	TCP  string `yaml:"tcp"`
	UDP  string `yaml:"udp"`
}

type BackendConfig struct {
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	Protocol string `yaml:"protocol"`
}

type DiscoveryConfig struct {
	Type        string `yaml:"type"`
	Address     string `yaml:"address"`
	ServiceName string `yaml:"service_name"`
	Protocol    string `yaml:"protocol"`
}

type PluginConfig struct {
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}

type Config struct {
	Listen    ListenConfig    `yaml:"listen"`
	Backends  []BackendConfig `yaml:"backends"`
	Discovery DiscoveryConfig `yaml:"discovery"`
	Plugins   []PluginConfig  `yaml:"plugins"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
