# devenv — Technical Specification

## 1. Architecture Overview
### Components
1) **devenv CLI** (host): user commands + repo MCP (`devenv stdio`)
2) **Host control plane (daemon)**: engine ops, worktrees, leases, instance state, per-instance control sockets
3) **In-container bridge**: internal MCP over stdio forwarding to the per-instance socket

Security principle: **no container engine socket in containers**; use a narrow, scoped RPC channel.

---

## 2. Data Models

### 2.1 Project config (`.devenv/config.yml`)
```yaml
version: 1

compose:
  template: .devenv/devenv-compose.yaml

agent:
  service: dev
  mount_path: /workspace
  workdir: /workspace
  default_shell: ["/bin/sh", "-lc"]

services:
  backend:  { compose_service: backend }
  frontend: { compose_service: frontend }
  db:       { compose_service: db }

service_groups:
  app: [backend, frontend]

ports:
  placeholders:
    FRONTEND_PORT: { range: "13000-13999", protocol: tcp }
    BACKEND_PORT:  { range: "13000-13999", protocol: tcp }

security:
  strict_port_exposure: true
```

### 2.2 Instance state (`.devenv/instances/<env>/state.json`)
```json
{
  "schema_version": 1,
  "repo_root": "/path/to/repo",
  "env_name": "feature-login",
  "instance_id": "uuid",
  "engine": "docker",
  "compose_project": "devenv_ab12_feature-login",
  "worktree_path": "/path/to/repo/.devenv/worktrees/feature-login",
  "compose_path": "/path/to/repo/.devenv/instances/feature-login/compose.yaml",
  "port_map": { "FRONTEND_PORT": 13027, "BACKEND_PORT": 13028 },
  "created_at": "2026-01-19T12:00:00Z",
  "agent": { "status": "running|stopped" },
  "capability": { "socket_path": "control.sock" }
}
```

### 2.3 Port lease DB (user-level, SQLite)
Tables:
- `leases(project_id, env_name, placeholder, host_port, created_at, last_used_at)`
- `ports(host_port PRIMARY KEY, reserved_by_project, reserved_by_env, placeholder, reserved_at)`

---

## 3. Engine/Compose Backend Abstraction
Interface:
- `up(project, compose_file, env_file?)`
- `down(project, compose_file, remove_volumes?)`
- `ps(project, compose_file)`
- `logs(project, compose_file, service, follow, tail)`
- `restart(project, compose_file, service)`
- `exec(project, compose_file, service, cmd, workdir?, env?)`

Backends:
- Docker: `docker compose`
- Podman: `podman compose` or docker-compat
- nerdctl: `nerdctl compose`

---

## 4. Template Rendering & Port Enforcement
- Placeholders: `{NAME}` where NAME is `[A-Z0-9_]+`
- Render via string substitution using leased values.
- Strict port exposure:
  - parse rendered YAML
  - ensure published host ports match leased placeholders/allowlist
  - otherwise fail `devenv up`

---

## 5. Worktree Management
- Create: `git worktree add .devenv/worktrees/<env> <ref>`
- Purge: `git worktree remove` + delete directory
- Policies: refuse if dirty unless `--force`; support `--reset/--recreate` later.

---

## 6. CLI Contract (v1)
### 6.1 Command Surface
- `devenv init`
- `devenv up <env>`
- `devenv down <env>`
- `devenv list`
- `devenv status <env>`
- `devenv ports <env>`
- `devenv service <op> <env> <svc> [-- ...]`
- `devenv agent <env> "<prompt>"`
- `devenv agent status <env>` / `devenv agent stop <env>`
- `devenv stdio`

### 6.2 Flags
- `init`: `--force`, `--template <path>`
- `up`: `--from <ref>`, `--agent "<prompt>"`, `--agent-file <path>`, `--no-agent`, `--engine <docker|podman|nerdctl>`, `--force`, `--ports <placeholder=port,...>`
- `down`: `--purge`, `--volumes`
- `list`: `--json`
- `status`: `--json`, `--verbose`
- `ports`: `--json`
- `service logs`: `--follow`, `--tail N`, `--since <duration>`, `--json`
- `service exec`: `--workdir <path>`, `--env KEY=VAL` (repeatable), `--user <uid:gid>`
- `agent`: `--agent <type>`, `--provider <p>`, `--model <m>`, `--prompt-file <path>`, `--prompt -`
- `stdio`: `--log-level`

### 6.3 Exit Codes
- `0` success
- `1` unexpected failure
- `2` invalid usage / bad flags / invalid env name
- `3` repo not initialized (`.devenv/config.yml` missing)
- `4` engine/compose backend failure
- `5` port lease failure / port conflicts
- `6` worktree conflict (dirty, exists, bad ref)
- `7` security violation (port exposure rule, scope)

### 6.4 Output Format
- Human-readable by default.
- JSON via `--json` for `list`, `status`, `ports`, and `service logs`.

### 6.5 Naming and Validation
- Normalize `<env>`: lowercase; spaces to `-`; collapse multiple `-`; trim leading/trailing `-`.
- Validate normalized `<env>` matches `^[a-z0-9][a-z0-9-_]{0,62}$`.
- Placeholders must match `{[A-Z0-9_]+}`.
- Service alias must exist in `.devenv/config.yml`.

---

## 7. Control Plane & Scoping
### 6.1 Per-instance control socket
- Create `.devenv/instances/<env>/control.sock`
- Generate capability token; store with restrictive permissions (or OS keychain)
- Mount into agent host container at `/run/devenv/control.sock`

### 6.2 RPC protocol
JSON messages over UDS:
- Request: `{id, method, params, token}`
- Response: `{id, ok, result}` or `{id, ok:false, error}`

Daemon enforces instance binding: socket+token map to exactly one compose project.

---

## 7. MCP Servers
### 7.1 Internal MCP (in-container): `devenv internal-mcp`
- stdio MCP server
- no `<env>` parameter; scope implicit
- tools:
  - `service.list` (list services)
  - `service.logs` (service logs)
  - `service.restart` (restart service)
  - `service.start` / `service.stop` (service lifecycle)
  - `ports.list` (leased ports)
  - `env.info` (env metadata)

### 7.2 Repo MCP (host): `devenv stdio`
- stdio MCP server scoped to current repo
- tools:
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

---

## 8. Agent Execution
- Start via `devenv up <env> --agent "<prompt>"` or `devenv agent <env> "<prompt>"`
- Prompt stored at `.devenv/instances/<env>/prompt.md`
- Agent runs inside agent host service using adapter-defined command.
- Agent integrates with internal MCP by launching `devenv internal-mcp` as an MCP tool.

---

## 9. Concurrency & Reliability
- Instance-level filesystem locks
- SQLite transactions for leasing
- `up` idempotent
- `status` reconciles stored state vs `compose ps`

---

## 10. Testing
- Unit: rendering, parsing, leasing, validators
- Integration: multi-env create/up/logs/restart/down
- Security: internal MCP cannot access other envs; no engine socket mounted
