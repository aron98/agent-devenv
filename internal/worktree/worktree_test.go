package worktree

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aron98/agent-devenv/internal/errs"
	"github.com/aron98/agent-devenv/internal/instance"
)

func TestCreateEnsureCleanAndRemove(t *testing.T) {
	repoRoot := setupRepo(t)
	layout := instance.Layout{RepoRoot: repoRoot}

	worktreePath, err := Create(layout, "feature-x", "HEAD")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if _, err := os.Stat(worktreePath); err != nil {
		t.Fatalf("worktree path missing: %v", err)
	}

	if err := EnsureClean(worktreePath, false); err != nil {
		t.Fatalf("EnsureClean returned error for clean worktree: %v", err)
	}

	dirtyFile := filepath.Join(worktreePath, "dirty.txt")
	if err := os.WriteFile(dirtyFile, []byte("dirty"), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	err = EnsureClean(worktreePath, false)
	if err == nil {
		t.Fatalf("EnsureClean returned nil for dirty worktree")
	}
	if errs.ExitCode(err) != exitWorktreeConflict {
		t.Fatalf("EnsureClean exit code = %d, want %d", errs.ExitCode(err), exitWorktreeConflict)
	}

	if err := Remove(layout, "feature-x"); err != nil {
		t.Fatalf("Remove returned error: %v", err)
	}
	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
		t.Fatalf("expected worktree path to be removed, got err %v", err)
	}
}

func TestPathUsesLayout(t *testing.T) {
	repoRoot := t.TempDir()
	layout := instance.Layout{RepoRoot: repoRoot}

	path, err := Path(layout, "alpha")
	if err != nil {
		t.Fatalf("Path returned error: %v", err)
	}
	if path != layout.WorktreeDir("alpha") {
		t.Fatalf("Path = %q, want %q", path, layout.WorktreeDir("alpha"))
	}
}

func setupRepo(t *testing.T) string {
	t.Helper()

	repoRoot := t.TempDir()
	runGitCommand(t, repoRoot, "init")
	runGitCommand(t, repoRoot, "config", "user.email", "test@example.com")
	runGitCommand(t, repoRoot, "config", "user.name", "Test User")

	readme := filepath.Join(repoRoot, "README.md")
	if err := os.WriteFile(readme, []byte("hello"), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	runGitCommand(t, repoRoot, "add", "README.md")
	runGitCommand(t, repoRoot, "commit", "-m", "init")

	return repoRoot
}

func runGitCommand(t *testing.T, dir string, args ...string) {
	t.Helper()

	if _, err := runGit(dir, args...); err != nil {
		t.Fatalf("git %v failed: %v", args, err)
	}
}
