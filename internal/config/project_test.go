package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

const validProjectConfigYAML = `version: 1
compose:
  template: .devenv/devenv-compose.yaml
agent:
  service: dev
  mount_path: /workspace
  workdir: /workspace
  default_shell: ["/bin/sh", "-lc"]
services:
  app:
    compose_service: app
ports:
  placeholders:
    FRONTEND_PORT:
      range: "13000-13999"
      protocol: tcp
security:
  strict_port_exposure: true
`

func TestLoadProjectConfig(t *testing.T) {
	path := writeTempConfig(t, validProjectConfigYAML)

	cfg, err := LoadProjectConfig(path)
	if err != nil {
		t.Fatalf("LoadProjectConfig returned error: %v", err)
	}

	if cfg.Version != 1 {
		t.Fatalf("Version = %d, want 1", cfg.Version)
	}
	if cfg.Compose.Template != ".devenv/devenv-compose.yaml" {
		t.Fatalf("Compose.Template = %q", cfg.Compose.Template)
	}
	if cfg.Agent.Service != "dev" {
		t.Fatalf("Agent.Service = %q", cfg.Agent.Service)
	}
	if cfg.Ports.Placeholders["FRONTEND_PORT"].Range != "13000-13999" {
		t.Fatalf("Ports.Placeholders.FRONTEND_PORT.Range = %q", cfg.Ports.Placeholders["FRONTEND_PORT"].Range)
	}
}

func TestValidateProjectConfig(t *testing.T) {
	path := writeTempConfig(t, validProjectConfigYAML)

	cfg, err := LoadProjectConfig(path)
	if err != nil {
		t.Fatalf("LoadProjectConfig returned error: %v", err)
	}
	if err := ValidateProjectConfig(cfg); err != nil {
		t.Fatalf("ValidateProjectConfig returned error: %v", err)
	}
}

func TestValidateProjectConfigMissingFields(t *testing.T) {
	var validationErr ValidationError

	err := ValidateProjectConfig(&ProjectConfig{})
	if err == nil {
		t.Fatalf("ValidateProjectConfig returned nil, want error")
	}
	if !errors.As(err, &validationErr) {
		t.Fatalf("ValidateProjectConfig error type = %T, want ValidationError", err)
	}
	if len(validationErr.Issues) == 0 {
		t.Fatalf("ValidationError.Issues = 0, want non-zero")
	}
}

func TestValidateProjectConfigNil(t *testing.T) {
	err := ValidateProjectConfig(nil)
	if err == nil {
		t.Fatalf("ValidateProjectConfig returned nil, want error")
	}
}

func writeTempConfig(t *testing.T, contents string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	return path
}
