package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const defaultLeaseRange = "13000-13999"

type UserConfig struct {
	Engines []string          `yaml:"engines"`
	Agent   UserAgentConfig   `yaml:"agent"`
	Secrets UserSecretsConfig `yaml:"secrets"`
	Ports   UserPortsConfig   `yaml:"ports"`
	Logging UserLoggingConfig `yaml:"logging"`
}

type UserAgentConfig struct {
	Type     string `yaml:"type"`
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
}

type UserSecretsConfig struct {
	Passthrough []string `yaml:"passthrough"`
}

type UserPortsConfig struct {
	LeaseDB      string `yaml:"lease_db"`
	DefaultRange string `yaml:"default_range"`
}

type UserLoggingConfig struct {
	Level string `yaml:"level"`
}

func UserConfigDir() (string, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}
	return filepath.Join(baseDir, "devenv"), nil
}

func UserConfigPath() (string, error) {
	configDir, err := UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yml"), nil
}

func DefaultUserConfig(configDir string) UserConfig {
	leaseDB := ""
	if configDir != "" {
		leaseDB = filepath.Join(configDir, "leases.db")
	}

	return UserConfig{
		Engines: []string{"docker", "podman", "nerdctl"},
		Ports: UserPortsConfig{
			LeaseDB:      leaseDB,
			DefaultRange: defaultLeaseRange,
		},
		Logging: UserLoggingConfig{
			Level: "info",
		},
	}
}

func LoadUserConfig(path string) (*UserConfig, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read user config: %w", err)
	}

	var cfg UserConfig
	if err := yaml.Unmarshal(contents, &cfg); err != nil {
		return nil, fmt.Errorf("parse user config: %w", err)
	}

	return &cfg, nil
}

func ResolveUserConfig(user *UserConfig, defaults UserConfig) UserConfig {
	resolved := defaults
	if user == nil {
		return resolved
	}

	if len(user.Engines) > 0 {
		resolved.Engines = append([]string{}, user.Engines...)
	}
	if user.Agent.Type != "" {
		resolved.Agent.Type = user.Agent.Type
	}
	if user.Agent.Provider != "" {
		resolved.Agent.Provider = user.Agent.Provider
	}
	if user.Agent.Model != "" {
		resolved.Agent.Model = user.Agent.Model
	}
	if len(user.Secrets.Passthrough) > 0 {
		resolved.Secrets.Passthrough = append([]string{}, user.Secrets.Passthrough...)
	}
	if user.Ports.LeaseDB != "" {
		resolved.Ports.LeaseDB = user.Ports.LeaseDB
	}
	if user.Ports.DefaultRange != "" {
		resolved.Ports.DefaultRange = user.Ports.DefaultRange
	}
	if user.Logging.Level != "" {
		resolved.Logging.Level = user.Logging.Level
	}

	return resolved
}
