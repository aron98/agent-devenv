package instance

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Layout struct {
	RepoRoot string
}

func (l Layout) BaseDir() string {
	return filepath.Join(l.RepoRoot, ".devenv")
}

func (l Layout) InstancesDir() string {
	return filepath.Join(l.BaseDir(), "instances")
}

func (l Layout) WorktreesDir() string {
	return filepath.Join(l.BaseDir(), "worktrees")
}

func (l Layout) InstanceDir(env string) string {
	return filepath.Join(l.InstancesDir(), env)
}

func (l Layout) WorktreeDir(env string) string {
	return filepath.Join(l.WorktreesDir(), env)
}

func (l Layout) ComposePath(env string) string {
	return filepath.Join(l.InstanceDir(env), "compose.yaml")
}

func (l Layout) InstanceEnvPath(env string) string {
	return filepath.Join(l.InstanceDir(env), "instance.env")
}

func (l Layout) StatePath(env string) string {
	return filepath.Join(l.InstanceDir(env), "state.json")
}

func (l Layout) ControlSocketPath(env string) string {
	return filepath.Join(l.InstanceDir(env), "control.sock")
}

func EnsureInstanceDirs(layout Layout, env string) error {
	if layout.RepoRoot == "" {
		return errors.New("repo root is required")
	}
	if env == "" {
		return errors.New("env name is required")
	}

	if err := os.MkdirAll(layout.InstancesDir(), 0o755); err != nil {
		return fmt.Errorf("create instances dir: %w", err)
	}
	if err := os.MkdirAll(layout.WorktreesDir(), 0o755); err != nil {
		return fmt.Errorf("create worktrees dir: %w", err)
	}
	if err := os.MkdirAll(layout.InstanceDir(env), 0o755); err != nil {
		return fmt.Errorf("create instance dir: %w", err)
	}

	return nil
}

func RemoveInstanceDir(layout Layout, env string) error {
	if layout.RepoRoot == "" {
		return errors.New("repo root is required")
	}
	if env == "" {
		return errors.New("env name is required")
	}

	if err := os.RemoveAll(layout.InstanceDir(env)); err != nil {
		return fmt.Errorf("remove instance dir: %w", err)
	}

	return nil
}
