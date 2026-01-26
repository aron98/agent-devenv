package worktree

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aron98/agent-devenv/internal/errs"
	"github.com/aron98/agent-devenv/internal/instance"
)

const exitWorktreeConflict = 6

func Path(layout instance.Layout, env string) (string, error) {
	if layout.RepoRoot == "" {
		return "", errors.New("repo root is required")
	}
	if env == "" {
		return "", errors.New("env name is required")
	}

	return layout.WorktreeDir(env), nil
}

func Create(layout instance.Layout, env string, ref string) (string, error) {
	if layout.RepoRoot == "" {
		return "", errors.New("repo root is required")
	}
	if env == "" {
		return "", errors.New("env name is required")
	}

	if err := os.MkdirAll(layout.WorktreesDir(), 0o755); err != nil {
		return "", fmt.Errorf("create worktrees dir: %w", err)
	}

	worktreePath := layout.WorktreeDir(env)
	args := []string{"worktree", "add", worktreePath}
	if strings.TrimSpace(ref) != "" {
		args = append(args, ref)
	}

	if _, err := runGit(layout.RepoRoot, args...); err != nil {
		return "", err
	}

	return worktreePath, nil
}

func EnsureClean(worktreePath string, force bool) error {
	if worktreePath == "" {
		return errors.New("worktree path is required")
	}

	clean, err := IsClean(worktreePath)
	if err != nil {
		return err
	}
	if clean || force {
		return nil
	}

	message := fmt.Sprintf("worktree at %s has uncommitted changes", worktreePath)
	return errs.New(exitWorktreeConflict, message, "use --force to proceed")
}

func IsClean(worktreePath string) (bool, error) {
	if worktreePath == "" {
		return false, errors.New("worktree path is required")
	}

	output, err := runGit(worktreePath, "status", "--porcelain")
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(string(output)) == "", nil
}

func Remove(layout instance.Layout, env string) error {
	if layout.RepoRoot == "" {
		return errors.New("repo root is required")
	}
	if env == "" {
		return errors.New("env name is required")
	}

	worktreePath := layout.WorktreeDir(env)
	_, err := runGit(layout.RepoRoot, "worktree", "remove", worktreePath)
	if err != nil {
		if removeErr := os.RemoveAll(worktreePath); removeErr != nil {
			return fmt.Errorf("remove worktree: %w", err)
		}
		return nil
	}

	if err := os.RemoveAll(worktreePath); err != nil {
		return fmt.Errorf("remove worktree dir: %w", err)
	}

	return nil
}

func runGit(dir string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}

	output, err := cmd.CombinedOutput()
	if err == nil {
		return output, nil
	}

	message := strings.TrimSpace(string(output))
	if message == "" {
		return output, fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}

	return output, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, message)
}
