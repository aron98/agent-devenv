# devenv

Go CLI for spinning up isolated, per-task development environments from a single repo. Each environment is a git worktree plus a dedicated Compose project with its own network, sidecars, and dynamically leased host ports. Agents run inside the environment container and are scoped to only their own services via an internal MCP bridge, keeping control channels isolated and safe.

Use `devenv init` to add a repo template, `devenv up <env>` to create and start an environment, and `devenv agent <env> "<prompt>"` to launch a coding agent inside it. The CLI also exposes a repo-scoped MCP server (`devenv stdio`) for orchestrators that need to list envs, manage services, and control agents without direct engine access.

## Development Status
- Phases 1-3 are complete: config parsing/validation, instance state/layout, and worktree management in internal packages.
- CLI commands are still stubbed while backend packages mature.
- CI runs `go build ./...` and `go test ./...` on PR open/sync and pushes to `main`.

## CLI Contract
- Commands: `init`, `up`, `down`, `list`, `status`, `ports`, `service <op>`, `agent`, `stdio`.
- `up` normalizes env names (lowercase, spaces to `-`, collapse `-`) and requires `--force` for dirty worktrees.
- Agent start supports inline prompt, `--agent-file`, or `--prompt -` (stdin).
- JSON output is supported for `list`, `status`, `ports`, and `service logs` via `--json`.
- Errors return a non-zero exit code with a short message and next step hint.

## MCP Capabilities
- **Repo MCP (`devenv stdio`)**: `env.list` (list envs), `env.up` (create/start env), `env.down` (stop/purge env), `env.status` (env status), `env.ports` (leased ports), `service.logs` (service logs), `service.restart` (restart service), `service.exec` (exec in service), `agent.start` (start agent), `agent.status` (agent status), `agent.stop` (stop agent).
- **Session MCP (`devenv internal-mcp`)**: `service.list` (list services), `service.logs` (service logs), `service.restart` (restart service), `service.start`/`service.stop` (service lifecycle), `ports.list` (leased ports), `env.info` (env metadata).
