# devenv — Product Requirements Document (PRD)

## Summary
**devenv** enables developers and orchestrators to create **any number of isolated, per-task development environments** from a single repository. Each environment is:
- A **git worktree** (separate checkout and working directory)
- A **Compose project** (containers + isolated network + optional sidecars)
- An **agent host container** where a coding agent runs **inside** the environment
- A **dynamic host port mapping** using placeholders (e.g., `{FRONTEND_PORT}`) leased per environment instance
- **Strictly scoped control**: agents can only read logs/control services **within their own environment**, via an internal MCP server and a per-environment control channel.

Users can start an environment and agent in one step:
- `devenv up <env> --agent "<prompt>"`

Or start/attach an agent later:
- `devenv agent <env> "<prompt>"`

A repo-scoped MCP server is available externally:
- `devenv stdio` (run in a repo with `.devenv/` initialized)

---

## Problem
Developers frequently need parallel work streams (feature vs bugfix vs spike) with repeatable local dependencies (DBs, queues, services), isolated changes, and consistent agent execution contexts. Today, teams use ad-hoc scripts, devcontainers, and Compose setups that are hard to clone, hard to keep isolated, and unsafe for agents (log/service control often requires access to the host container engine socket).

---

## Goals
1. **Unlimited environments per repo**: users create environments on demand with arbitrary names.
2. **Isolation**:
   - Separate worktrees (no file conflicts)
   - Separate Compose projects/networks (no port conflicts)
3. **Agent runs in-env**:
   - Agent executes within an environment’s dev container
   - Can access sidecars via the environment network
4. **Scoped operations**:
   - In-env agent can only control services/logs of its own env (not other envs)
5. **Dynamic ports**:
   - Projects declare placeholders `{NAME_PORT}` and devenv leases host ports per env
6. **Engine flexibility**:
   - Works with Docker + alternatives (Podman, nerdctl)
7. **Orchestrator integration**:
   - Repo-scoped MCP server with environment lifecycle, logs, service control, agent control

---

## Non-Goals (v1)
- Kubernetes-native orchestration (Compose only)
- Multi-repo / monorepo orchestration beyond a single repo root (possible later)
- Built-in secrets manager (v1 supports env passthrough and local `.env` patterns)
- Remote execution / hosted environments
- Full devcontainer spec implementation beyond “Compose + mounted worktree + dev service”

---

## Personas & Use Cases
### Individual Developer
- Spin up `feature-login` and `bug-217` simultaneously
- Run separate DB state per environment (or shared volume if configured)
- Start an agent in each env with different prompts

### Team Lead / Maintainer
- Commit `.devenv/` project templates so everyone has consistent env topology
- Ensure ports exposed are only those intended

### CI / Orchestrator
- Programmatically create envs, run commands, stream logs, and teardown via MCP

---

## Key User Stories
1. As a developer, I can run `devenv init` to add a minimal, shareable project template.
2. As a developer, I can run `devenv up feature-x` to create a worktree + start sidecars in an isolated network.
3. As a developer, I can run `devenv up feature-x --agent "…"`, which also starts an agent inside the env.
4. As a developer, I can run `devenv service logs feature-x backend` to view backend logs in that env.
5. As a developer, I can run `devenv service restart feature-x backend` to restart only that service in that env.
6. As an agent running inside an env, I can read logs and restart services only in that env (never across envs).
7. As an orchestrator, I can run `devenv stdio` and control only envs belonging to the repo.

---

## User Experience Overview
### Project setup
- `devenv init` creates:
  - `.devenv/devenv-compose.yaml` (template with placeholders)
  - `.devenv/config.yml` (service aliases, port placeholders, agent host service)
- It also ensures generated runtime directories are gitignored:
  - `.devenv/worktrees/`
  - `.devenv/instances/`

### Environment lifecycle
- `devenv up <env> [--from <ref>]`:
  1) normalize `<env>` (lowercase; spaces to `-`; collapse `-`)
  2) create worktree or require clean unless `--force`
  3) lease ports (overrideable with `--ports`)
  4) render compose template
  5) start compose project
  6) optionally start agent in-env
- `devenv down <env> [--purge] [--volumes]`:
  - stop project
  - `--purge` removes worktree + releases leases + deletes state

### Agent lifecycle
- `devenv agent <env> "<prompt>"`:
  - ensures env is running
  - starts or attaches agent inside agent host service
  - supports `--prompt-file` or `--prompt -` (stdin)
- `devenv agent status <env>` and `devenv agent stop <env>` (v1 recommended)

### Service lifecycle
- `devenv service logs <env> <svc> [--follow] [--tail N] [--since <duration>]`
- `devenv service restart <env> <svc>`
- `devenv service exec <env> <svc> -- <cmd...>`

### Repo MCP server
- `devenv stdio` exposes MCP tools scoped to this repo’s `.devenv/instances/*`

---
## MCP Capabilities (v1)
### Repo MCP (`devenv stdio`)
- `env.list` (list envs)
- `env.up` (create/start env)
- `env.down` (stop/purge env)
- `env.status` (env status)
- `env.ports` (leased ports)
- `service.logs` (service logs)
- `service.restart` (restart service)
- `service.exec` (exec in service)
- `agent.start` (start agent)
- `agent.status` (agent status)
- `agent.stop` (stop agent)

### Session MCP (`devenv internal-mcp`)
- `service.list` (list services)
- `service.logs` (service logs)
- `service.restart` (restart service)
- `service.start` / `service.stop` (service lifecycle)
- `ports.list` (leased ports)
- `env.info` (env metadata)

---
## CLI Contract (v1)
### Commands
- `init`, `up`, `down`, `list`, `status`, `ports`, `service <op>`, `agent`, `stdio`.

### Flags
- `init`: `--force`, `--template <path>`
- `up`: `--from <ref>`, `--agent "<prompt>"`, `--agent-file <path>`, `--no-agent`, `--engine <docker|podman|nerdctl>`, `--force`, `--ports <placeholder=port,...>`
- `down`: `--purge`, `--volumes`
- `list`: `--json`
- `status`: `--json`, `--verbose`
- `ports`: `--json`
- `service logs`: `--follow`, `--tail N`, `--since <duration>`, `--json`
- `service exec`: `--workdir <path>`, `--env KEY=VAL`, `--user <uid:gid>`
- `agent`: `--agent <type>`, `--provider <p>`, `--model <m>`, `--prompt-file <path>`, `--prompt -`
- `stdio`: `--log-level`

### Output
- Human-readable by default.
- JSON is supported for `list`, `status`, `ports`, and `service logs` via `--json`.

### Errors
- Non-zero exit code with short message and next-step hint.

### Naming
- Normalize `<env>`: lowercase; spaces to `-`; collapse multiple `-`; trim `-`.
- Validate normalized `<env>` matches `^[a-z0-9][a-z0-9-_]{0,62}$`.

---

## Requirements

### Functional Requirements
- FR1: Create arbitrary environment instances without predeclared env config.
- FR2: Each env has an isolated worktree and isolated Compose project/network.
- FR3: Template compose supports placeholders for host ports (e.g., `{FRONTEND_PORT}`).
- FR4: Port leases are stable across restarts for the same env.
- FR5: Agent runs inside env in a dedicated Compose service (agent host).
- FR6: Agent can read logs/control services only in its env via internal MCP.
- FR7: Host CLI can list envs, show status, logs, restart services, and show ports.
- FR8: Works with Docker, Podman, and nerdctl backends (best-effort within Compose compatibility).
- FR9: Repo-scoped MCP server (`devenv stdio`) controls env lifecycle, logs, services, agent.

### Non-Functional Requirements
- NFR1: Safe-by-default — no container engine socket mounted into agent containers.
- NFR2: Strong scoping — in-env control channels cannot operate outside the instance.
- NFR3: Good ergonomics — commands are consistent and shell-completable.
- NFR4: Deterministic file layout and recoverability after crashes.
- NFR5: Concurrency-safe — parallel `devenv up` for different envs works reliably.
- NFR6: Cross-platform where feasible (Linux/macOS first; Windows later or best-effort).

---

## Success Metrics
- Time to first working environment (from clone + `devenv init`) under 5 minutes for typical repos.
- 90%+ of common service operations (logs/restart/ports) require no manual Compose commands.
- Agent safety: no cross-env control possible via internal MCP.
- Low friction multi-env: two or more envs can run simultaneously with no port collisions.

---

## Risks & Mitigations
- **Compose feature differences** across docker/podman/nerdctl  
  Mitigation: validate template subset; provide clear errors and compatibility docs.
- **Port leasing collisions**  
  Mitigation: single lease DB with locking; preflight port availability; stable reservations.
- **Agent diversity** (opencode/goose/claude-code/codex)  
  Mitigation: “adapter” approach; default to running a shell command inside agent host.
- **State drift** after crashes  
  Mitigation: state.json + reconcile commands; `devenv doctor` (post-v1) and robust `status`.

---

## Out of Scope Decisions (to revisit)
- Multi-user shared machines
- Remote/cluster execution
- Multi-repo “workspace sets”
- Built-in secret storage
