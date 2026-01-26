package instance

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

const SchemaVersion = 1

type AgentState struct {
	Status string `json:"status"`
}

type CapabilityState struct {
	SocketPath string `json:"socket_path"`
	Token      string `json:"token,omitempty"`
}

type InstanceState struct {
	SchemaVersion  int             `json:"schema_version"`
	RepoRoot       string          `json:"repo_root"`
	EnvName        string          `json:"env_name"`
	InstanceID     string          `json:"instance_id"`
	Engine         string          `json:"engine"`
	ComposeProject string          `json:"compose_project"`
	WorktreePath   string          `json:"worktree_path"`
	ComposePath    string          `json:"compose_path"`
	PortMap        map[string]int  `json:"port_map"`
	CreatedAt      time.Time       `json:"created_at"`
	Agent          AgentState      `json:"agent"`
	Capability     CapabilityState `json:"capability"`
}

func WriteState(layout Layout, env string, state InstanceState) error {
	if env == "" {
		return errors.New("env name is required")
	}
	if state.SchemaVersion == 0 {
		state.SchemaVersion = SchemaVersion
	}
	if state.EnvName == "" {
		state.EnvName = env
	}
	if err := EnsureInstanceDirs(layout, env); err != nil {
		return err
	}

	payload, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}

	if err := os.WriteFile(layout.StatePath(env), payload, 0o600); err != nil {
		return fmt.Errorf("write state: %w", err)
	}

	return nil
}

func LoadState(layout Layout, env string) (InstanceState, error) {
	if env == "" {
		return InstanceState{}, errors.New("env name is required")
	}

	payload, err := os.ReadFile(layout.StatePath(env))
	if err != nil {
		return InstanceState{}, fmt.Errorf("read state: %w", err)
	}

	var state InstanceState
	if err := json.Unmarshal(payload, &state); err != nil {
		return InstanceState{}, fmt.Errorf("parse state: %w", err)
	}

	return state, nil
}
