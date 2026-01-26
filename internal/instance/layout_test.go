package instance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureInstanceDirsCreatesPaths(t *testing.T) {
	repoRoot := t.TempDir()
	layout := Layout{RepoRoot: repoRoot}

	if err := EnsureInstanceDirs(layout, "feature-x"); err != nil {
		t.Fatalf("EnsureInstanceDirs returned error: %v", err)
	}

	paths := []string{
		layout.InstancesDir(),
		layout.WorktreesDir(),
		layout.InstanceDir("feature-x"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}
}

func TestRemoveInstanceDirRemovesOnlyInstance(t *testing.T) {
	repoRoot := t.TempDir()
	layout := Layout{RepoRoot: repoRoot}

	if err := EnsureInstanceDirs(layout, "alpha"); err != nil {
		t.Fatalf("EnsureInstanceDirs returned error: %v", err)
	}

	if err := RemoveInstanceDir(layout, "alpha"); err != nil {
		t.Fatalf("RemoveInstanceDir returned error: %v", err)
	}

	if _, err := os.Stat(layout.InstanceDir("alpha")); !os.IsNotExist(err) {
		t.Fatalf("expected instance dir to be removed, got err %v", err)
	}

	if _, err := os.Stat(filepath.Join(repoRoot, ".devenv", "instances")); err != nil {
		t.Fatalf("expected instances dir to exist: %v", err)
	}
}
