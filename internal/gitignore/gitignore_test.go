package gitignore

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureEntriesCreatesFileAndIsIdempotent(t *testing.T) {
	repoRoot := t.TempDir()
	path := filepath.Join(repoRoot, ".gitignore")
	entries := []string{".devenv/worktrees/", ".devenv/instances/"}

	updated, err := EnsureEntries(path, entries)
	if err != nil {
		t.Fatalf("EnsureEntries returned error: %v", err)
	}
	if !updated {
		t.Fatalf("EnsureEntries updated = false, want true")
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	text := string(contents)
	for _, entry := range entries {
		if !strings.Contains(text, entry) {
			t.Fatalf("expected entry %q in gitignore", entry)
		}
	}

	updated, err = EnsureEntries(path, entries)
	if err != nil {
		t.Fatalf("EnsureEntries returned error on second run: %v", err)
	}
	if updated {
		t.Fatalf("EnsureEntries updated = true on second run, want false")
	}

	contentsAfter, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error after second run: %v", err)
	}
	if string(contentsAfter) != text {
		t.Fatalf("gitignore contents changed on second run")
	}
}
