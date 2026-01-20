package config

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Issues []string
}

func (err ValidationError) Error() string {
	return fmt.Sprintf("config validation failed: %s", err.Issues)
}

func (err ValidationError) Is(target error) bool {
	_, ok := target.(ValidationError)
	return ok
}

func ValidateProjectConfig(cfg *ProjectConfig) error {
	if cfg == nil {
		return errors.New("config is required")
	}

	var issues []string
	if cfg.Version == 0 {
		issues = append(issues, "version is required")
	}
	if cfg.Compose.Template == "" {
		issues = append(issues, "compose.template is required")
	}
	if cfg.Agent.Service == "" {
		issues = append(issues, "agent.service is required")
	}
	if cfg.Agent.MountPath == "" {
		issues = append(issues, "agent.mount_path is required")
	}
	if cfg.Agent.Workdir == "" {
		issues = append(issues, "agent.workdir is required")
	}
	if len(cfg.Agent.DefaultShell) == 0 {
		issues = append(issues, "agent.default_shell is required")
	}
	if len(cfg.Services) == 0 {
		issues = append(issues, "services is required")
	}
	if len(cfg.Ports.Placeholders) == 0 {
		issues = append(issues, "ports.placeholders is required")
	}

	if len(issues) > 0 {
		return ValidationError{Issues: issues}
	}

	return nil
}
