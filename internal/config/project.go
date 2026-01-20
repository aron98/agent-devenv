package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ProjectConfig struct {
	Version       int                    `yaml:"version"`
	Compose       ComposeConfig           `yaml:"compose"`
	Agent         AgentConfig             `yaml:"agent"`
	Services      map[string]ServiceAlias `yaml:"services"`
	ServiceGroups map[string][]string     `yaml:"service_groups"`
	Ports         PortsConfig             `yaml:"ports"`
	Security      SecurityConfig          `yaml:"security"`
}

type ComposeConfig struct {
	Template string `yaml:"template"`
}

type AgentConfig struct {
	Service      string   `yaml:"service"`
	MountPath    string   `yaml:"mount_path"`
	Workdir      string   `yaml:"workdir"`
	DefaultShell []string `yaml:"default_shell"`
}

type ServiceAlias struct {
	ComposeService string `yaml:"compose_service"`
}

type PortsConfig struct {
	Placeholders map[string]PortPlaceholder `yaml:"placeholders"`
}

type PortPlaceholder struct {
	Range    string `yaml:"range"`
	Protocol string `yaml:"protocol"`
}

type SecurityConfig struct {
	StrictPortExposure bool `yaml:"strict_port_exposure"`
}

func LoadProjectConfig(path string) (*ProjectConfig, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read project config: %w", err)
	}

	var cfg ProjectConfig
	if err := yaml.Unmarshal(contents, &cfg); err != nil {
		return nil, fmt.Errorf("parse project config: %w", err)
	}

	return &cfg, nil
}
