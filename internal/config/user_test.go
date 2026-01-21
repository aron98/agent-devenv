package config

import (
	"os"
	"path/filepath"
	"testing"
)

const userConfigYAML = `engines:
  - podman
agent:
  type: opencode
  provider: openai
  model: gpt-4
secrets:
  passthrough:
    - AWS_PROFILE
ports:
  lease_db: /tmp/leases.db
  default_range: "14000-14999"
logging:
  level: debug
`

func TestDefaultUserConfig(t *testing.T) {
	baseDir := filepath.Join(t.TempDir(), "devenv")
	cfg := DefaultUserConfig(baseDir)

	if cfg.Ports.LeaseDB != filepath.Join(baseDir, "leases.db") {
		t.Fatalf("LeaseDB = %q", cfg.Ports.LeaseDB)
	}
	if cfg.Ports.DefaultRange != defaultLeaseRange {
		t.Fatalf("DefaultRange = %q", cfg.Ports.DefaultRange)
	}
	if cfg.Logging.Level != "info" {
		t.Fatalf("Logging.Level = %q", cfg.Logging.Level)
	}
	if len(cfg.Engines) == 0 {
		t.Fatalf("Engines should not be empty")
	}
}

func TestLoadUserConfigMissingFile(t *testing.T) {
	cfg, err := LoadUserConfig(filepath.Join(t.TempDir(), "missing.yml"))
	if err != nil {
		t.Fatalf("LoadUserConfig returned error: %v", err)
	}
	if cfg != nil {
		t.Fatalf("LoadUserConfig returned config, want nil")
	}
}

func TestLoadUserConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(path, []byte(userConfigYAML), 0o600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	cfg, err := LoadUserConfig(path)
	if err != nil {
		t.Fatalf("LoadUserConfig returned error: %v", err)
	}
	if cfg.Agent.Type != "opencode" {
		t.Fatalf("Agent.Type = %q", cfg.Agent.Type)
	}
	if cfg.Ports.DefaultRange != "14000-14999" {
		t.Fatalf("Ports.DefaultRange = %q", cfg.Ports.DefaultRange)
	}
}

func TestResolveUserConfig(t *testing.T) {
	defaults := DefaultUserConfig("/tmp/devenv")
	user := &UserConfig{
		Engines: []string{"podman"},
		Ports: UserPortsConfig{
			DefaultRange: "14000-14999",
		},
		Logging: UserLoggingConfig{
			Level: "debug",
		},
	}

	resolved := ResolveUserConfig(user, defaults)
	if len(resolved.Engines) != 1 || resolved.Engines[0] != "podman" {
		t.Fatalf("Engines = %v", resolved.Engines)
	}
	if resolved.Ports.DefaultRange != "14000-14999" {
		t.Fatalf("DefaultRange = %q", resolved.Ports.DefaultRange)
	}
	if resolved.Ports.LeaseDB == "" {
		t.Fatalf("LeaseDB should be preserved from defaults")
	}
	if resolved.Logging.Level != "debug" {
		t.Fatalf("Logging.Level = %q", resolved.Logging.Level)
	}
}
