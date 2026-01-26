# Implementation Plan

A phased checklist for delivering a fully functional devenv CLI and MCP stack. Each step is intended to be small and verifiable.

## Phase 0 - Project scaffolding
- [x] Initialize Go module and basic repo structure. (test: `go test ./...`)
- [x] Add minimal CLI entrypoint with command routing. (test: `go run ./cmd/devenv help`)
- [x] Define core packages (config, worktree, ports, compose, control, mcp). (test: `go test ./internal/...`)
- [x] Establish logging and error helpers. (test: `go test ./internal/errs ./internal/logging`)

## Phase 1 - Config and validation
- [x] Parse `.devenv/config.yml` into typed config structs. (test: unit tests for config loading)
- [x] Validate config schema and required fields. (test: unit tests for `ValidateProjectConfig`)
- [x] Implement env name normalization (lowercase, spaces to `-`, collapse `-`, trim `-`). (test: unit tests for `NormalizeEnvName`)
- [x] Validate normalized env name against `^[a-z0-9][a-z0-9-_]{0,62}$`. (test: unit tests for `ValidateEnvName`)
- [x] Load user config and resolve precedence. (test: unit tests for config precedence rules)

## Phase 2 - Instance state and filesystem layout
- [x] Define instance state schema and structs. (test: unit tests for state serialization)
- [x] Implement instance directory creation and cleanup. (test: unit tests with temp directories)
- [x] Write and read `.devenv/instances/<env>/state.json`. (test: round-trip state tests)
- [x] Ensure `.devenv/worktrees/` and `.devenv/instances/` are gitignored. (test: init output asserts `.gitignore` entries)

## Phase 3 - Worktree management
- [ ] Implement worktree create (`git worktree add`) with ref support. (test: integration test using temp repo)
- [ ] Detect dirty worktrees and require `--force` to proceed. (test: integration test with modified file)
- [ ] Implement worktree removal for purge. (test: integration test removing worktree)
- [ ] Add unit tests for worktree path and normalization behavior. (test: `go test ./internal/worktree`)

## Phase 4 - Port leasing
- [ ] Create SQLite lease DB schema and migrations. (test: migration tests with temp DB)
- [ ] Implement lease acquisition for placeholders. (test: unit tests for placeholder leasing)
- [ ] Implement lease release on purge. (test: unit tests for release logic)
- [ ] Enforce stable leases per env across restarts. (test: repeat lease acquisition tests)
- [ ] Add override support for `--ports`. (test: unit tests for overrides)
- [ ] Validate leased ports are available before use. (test: unit tests with occupied port detection)

## Phase 5 - Compose template rendering
- [ ] Parse template and substitute `{PLACEHOLDER}` values. (test: unit tests for rendering)
- [ ] Write rendered compose to `.devenv/instances/<env>/compose.yaml`. (test: temp dir write tests)
- [ ] Enforce strict port exposure rules from config. (test: unit tests for enforcement)
- [ ] Support optional `.devenv/instances/<env>/instance.env`. (test: unit tests for env file creation)

## Phase 6 - Compose backend
- [ ] Implement backend interface for `up`, `down`, `ps`, `logs`, `restart`, `exec`. (test: interface unit tests with fakes)
- [ ] Add Docker Compose backend. (test: integration test with docker compose)
- [ ] Add Podman Compose backend. (test: integration test with podman compose)
- [ ] Add nerdctl Compose backend. (test: integration test with nerdctl compose)
- [ ] Normalize backend errors to CLI exit codes. (test: unit tests for exit code mapping)

## Phase 7 - Control plane sockets
- [ ] Create per-instance control socket and capability token. (test: unit tests for token creation)
- [ ] Persist socket path and token in instance state. (test: state write/read tests)
- [ ] Mount socket into agent host container. (test: integration test verifying mount)
- [ ] Enforce env scoping in control plane. (test: security tests for cross-env denial)

## Phase 8 - MCP servers
- [ ] Implement internal MCP (`devenv internal-mcp`). (test: unit tests for server startup)
- [ ] Add tools: `service.list`, `service.logs`, `service.restart`, `service.start`, `service.stop`, `ports.list`, `env.info`. (test: tool-level unit tests)
- [ ] Implement repo MCP (`devenv stdio`). (test: unit tests for server startup)
- [ ] Add tools: `env.list`, `env.up`, `env.down`, `env.status`, `env.ports`, `service.logs`, `service.restart`, `service.exec`, `agent.start`, `agent.status`, `agent.stop`. (test: tool-level unit tests)
- [ ] Validate internal MCP cannot control other envs. (test: security tests for scope)

## Phase 9 - Agent execution
- [ ] Store agent prompt at `.devenv/instances/<env>/prompt.md`. (test: unit tests for prompt persistence)
- [ ] Implement agent start inside agent host service. (test: integration test for agent start)
- [ ] Implement `devenv agent status` and `devenv agent stop`. (test: integration test for lifecycle)
- [ ] Ensure internal MCP is available to agents. (test: integration test verifying tool availability)

## Phase 10 - CLI commands
- [ ] Implement `devenv init`. (test: integration test for generated files)
- [ ] Implement `devenv up` including worktree, ports, render, compose up. (test: integration test for env creation)
- [ ] Implement `devenv down` with `--purge` and `--volumes`. (test: integration test for teardown)
- [ ] Implement `devenv list`. (test: integration test for list output)
- [ ] Implement `devenv status`. (test: integration test for status output)
- [ ] Implement `devenv ports`. (test: integration test for ports output)
- [ ] Implement `devenv service logs|restart|exec`. (test: integration tests for service ops)
- [ ] Implement `devenv stdio`. (test: integration test for MCP server)
- [ ] Wire exit codes and error hints. (test: unit tests for error mapping)

## Phase 11 - Output formats
- [ ] Add JSON output for `list`, `status`, `ports`, and `service logs`. (test: unit tests for JSON output)
- [ ] Ensure stable JSON fields per CLI contract. (test: golden JSON tests)

## Phase 12 - Testing and validation
- [ ] Unit tests for normalization, rendering, and leasing. (test: `go test ./internal/...`)
- [ ] Integration tests for `up/down/list/status/ports`. (test: integration suite)
- [ ] Integration tests for service logs/restart/exec. (test: integration suite)
- [ ] Security tests for MCP scoping and port exposure enforcement. (test: security suite)

## Phase 13 - Release readiness
- [ ] Add versioning and build metadata. (test: verify version flags)
- [ ] Create smoke test checklist for a sample repo. (test: manual smoke run)
- [ ] Verify docs match CLI behavior and MCP capabilities. (test: doc review checklist)
