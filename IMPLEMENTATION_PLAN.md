# Implementation Plan

A phased checklist for delivering a fully functional devenv CLI and MCP stack. Each step is intended to be small and verifiable.

## Phase 0 - Project scaffolding
- [x] Initialize Go module and basic repo structure.
- [ ] Add minimal CLI entrypoint with command routing.
- [ ] Define core packages (config, worktree, ports, compose, control, mcp).
- [ ] Establish logging and error helpers.

## Phase 1 - Config and validation
- [x] Parse `.devenv/config.yml` into typed config structs.
- [x] Validate config schema and required fields.
- [x] Implement env name normalization (lowercase, spaces to `-`, collapse `-`, trim `-`).
- [x] Validate normalized env name against `^[a-z0-9][a-z0-9-_]{0,62}$`.
- [ ] Load user config and resolve precedence.

## Phase 2 - Instance state and filesystem layout
- [ ] Define instance state schema and structs.
- [ ] Implement instance directory creation and cleanup.
- [ ] Write and read `.devenv/instances/<env>/state.json`.
- [ ] Ensure `.devenv/worktrees/` and `.devenv/instances/` are gitignored.

## Phase 3 - Worktree management
- [ ] Implement worktree create (`git worktree add`) with ref support.
- [ ] Detect dirty worktrees and require `--force` to proceed.
- [ ] Implement worktree removal for purge.
- [ ] Add unit tests for worktree path and normalization behavior.

## Phase 4 - Port leasing
- [ ] Create SQLite lease DB schema and migrations.
- [ ] Implement lease acquisition for placeholders.
- [ ] Implement lease release on purge.
- [ ] Enforce stable leases per env across restarts.
- [ ] Add override support for `--ports`.
- [ ] Validate leased ports are available before use.

## Phase 5 - Compose template rendering
- [ ] Parse template and substitute `{PLACEHOLDER}` values.
- [ ] Write rendered compose to `.devenv/instances/<env>/compose.yaml`.
- [ ] Enforce strict port exposure rules from config.
- [ ] Support optional `.devenv/instances/<env>/instance.env`.

## Phase 6 - Compose backend
- [ ] Implement backend interface for `up`, `down`, `ps`, `logs`, `restart`, `exec`.
- [ ] Add Docker Compose backend.
- [ ] Add Podman Compose backend.
- [ ] Add nerdctl Compose backend.
- [ ] Normalize backend errors to CLI exit codes.

## Phase 7 - Control plane sockets
- [ ] Create per-instance control socket and capability token.
- [ ] Persist socket path and token in instance state.
- [ ] Mount socket into agent host container.
- [ ] Enforce env scoping in control plane.

## Phase 8 - MCP servers
- [ ] Implement internal MCP (`devenv internal-mcp`).
- [ ] Add tools: `service.list`, `service.logs`, `service.restart`, `service.start`, `service.stop`, `ports.list`, `env.info`.
- [ ] Implement repo MCP (`devenv stdio`).
- [ ] Add tools: `env.list`, `env.up`, `env.down`, `env.status`, `env.ports`, `service.logs`, `service.restart`, `service.exec`, `agent.start`, `agent.status`, `agent.stop`.
- [ ] Validate internal MCP cannot control other envs.

## Phase 9 - Agent execution
- [ ] Store agent prompt at `.devenv/instances/<env>/prompt.md`.
- [ ] Implement agent start inside agent host service.
- [ ] Implement `devenv agent status` and `devenv agent stop`.
- [ ] Ensure internal MCP is available to agents.

## Phase 10 - CLI commands
- [ ] Implement `devenv init`.
- [ ] Implement `devenv up` including worktree, ports, render, compose up.
- [ ] Implement `devenv down` with `--purge` and `--volumes`.
- [ ] Implement `devenv list`.
- [ ] Implement `devenv status`.
- [ ] Implement `devenv ports`.
- [ ] Implement `devenv service logs|restart|exec`.
- [ ] Implement `devenv stdio`.
- [ ] Wire exit codes and error hints.

## Phase 11 - Output formats
- [ ] Add JSON output for `list`, `status`, `ports`, and `service logs`.
- [ ] Ensure stable JSON fields per CLI contract.

## Phase 12 - Testing and validation
- [ ] Unit tests for normalization, rendering, and leasing.
- [ ] Integration tests for `up/down/list/status/ports`.
- [ ] Integration tests for service logs/restart/exec.
- [ ] Security tests for MCP scoping and port exposure enforcement.

## Phase 13 - Release readiness
- [ ] Add versioning and build metadata.
- [ ] Create smoke test checklist for a sample repo.
- [ ] Verify docs match CLI behavior and MCP capabilities.
