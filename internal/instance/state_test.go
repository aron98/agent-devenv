package instance

import (
	"path/filepath"
	"testing"
	"time"
)

func TestStateRoundTrip(t *testing.T) {
	repoRoot := t.TempDir()
	layout := Layout{RepoRoot: repoRoot}

	createdAt := time.Date(2026, 1, 19, 12, 0, 0, 0, time.UTC)
	state := InstanceState{
		RepoRoot:       repoRoot,
		EnvName:        "feature-login",
		InstanceID:     "instance-123",
		Engine:         "docker",
		ComposeProject: "devenv_ab12_feature-login",
		WorktreePath:   filepath.Join(repoRoot, ".devenv", "worktrees", "feature-login"),
		ComposePath:    filepath.Join(repoRoot, ".devenv", "instances", "feature-login", "compose.yaml"),
		PortMap: map[string]int{
			"FRONTEND_PORT": 13027,
			"BACKEND_PORT":  13028,
		},
		CreatedAt:  createdAt,
		Agent:      AgentState{Status: "running"},
		Capability: CapabilityState{SocketPath: "control.sock"},
	}

	if err := WriteState(layout, "feature-login", state); err != nil {
		t.Fatalf("WriteState returned error: %v", err)
	}

	loaded, err := LoadState(layout, "feature-login")
	if err != nil {
		t.Fatalf("LoadState returned error: %v", err)
	}

	if loaded.SchemaVersion != SchemaVersion {
		t.Fatalf("SchemaVersion = %d, want %d", loaded.SchemaVersion, SchemaVersion)
	}
	if loaded.EnvName != state.EnvName {
		t.Fatalf("EnvName = %q, want %q", loaded.EnvName, state.EnvName)
	}
	if loaded.ComposeProject != state.ComposeProject {
		t.Fatalf("ComposeProject = %q, want %q", loaded.ComposeProject, state.ComposeProject)
	}
	if loaded.WorktreePath != state.WorktreePath {
		t.Fatalf("WorktreePath = %q, want %q", loaded.WorktreePath, state.WorktreePath)
	}
	if loaded.ComposePath != state.ComposePath {
		t.Fatalf("ComposePath = %q, want %q", loaded.ComposePath, state.ComposePath)
	}
	if loaded.PortMap["FRONTEND_PORT"] != 13027 {
		t.Fatalf("PortMap.FRONTEND_PORT = %d", loaded.PortMap["FRONTEND_PORT"])
	}
	if !loaded.CreatedAt.Equal(createdAt) {
		t.Fatalf("CreatedAt = %v, want %v", loaded.CreatedAt, createdAt)
	}
	if loaded.Agent.Status != state.Agent.Status {
		t.Fatalf("Agent.Status = %q, want %q", loaded.Agent.Status, state.Agent.Status)
	}
	if loaded.Capability.SocketPath != state.Capability.SocketPath {
		t.Fatalf("Capability.SocketPath = %q, want %q", loaded.Capability.SocketPath, state.Capability.SocketPath)
	}
}
