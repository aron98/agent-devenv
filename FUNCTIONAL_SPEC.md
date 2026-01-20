# devenv — Functional Specification

## 1. Concepts & Terminology
- **Repo**: a directory containing `.devenv/` initialized files.
- **Environment instance** (`<env>`): an arbitrary user-chosen name representing one isolated dev environment.
- **Worktree**: git worktree directory created per environment instance.
- **Compose project**: container stack started per environment instance (isolated network, containers, volumes as declared).
- **Agent host service**: a Compose service (e.g., `dev`) where the coding agent runs inside the environment.
- **Sidecars**: additional Compose services (db/redis/etc.) in the same network.
- **Port placeholders**: tokens in `.devenv/devenv-compose.yaml` like `{FRONTEND_PORT}` to be replaced by leased host ports.
- **Repo MCP server**: `devenv stdio` exposes repo-scoped control for orchestrators.
- **Internal MCP server**: runs inside the agent host container; only controls that environment.

---

## 2. Configuration Model

### 2.1 Project-level configuration (committed)
**Directory**: `.devenv/`

#### `.devenv/devenv-compose.yaml` (template)
- Compose YAML template with placeholders (e.g., `{FRONTEND_PORT}`).
- Must include an agent host service (configurable name via `.devenv/config.yml`).
- Services should only expose host ports using placeholders defined in `.devenv/config.yml`.

#### `.devenv/config.yml`
Defines:
- compose template path
- agent host service name, mount path, default workdir
- service aliases -> compose service names
- optional service groups
- port placeholder names and ranges/policies
- strict rules for port exposure
- optional health checks and bootstrap commands

**Important**: This file MUST NOT define environment instances. Instances are user-created.

### 2.2 User-level configuration (local, not committed)
Stored under OS-appropriate config directory (e.g., `~/.config/devenv/config.yml`).
Defines:
- engine preference order (docker/podman/nerdctl)
- default agent type + provider/model defaults
- secret env passthrough allowlist
- port lease DB location and default lease range
- logging verbosity

### 2.3 Precedence
1. CLI flags
2. Environment instance state (rendered files / state.json)
3. Project config (`.devenv/config.yml`)
4. User config
5. Built-in defaults

---

## 3. Filesystem Layout

### 3.1 Repo-managed (committed)
```
.devenv/
  devenv-compose.yaml
  config.yml
```

### 3.2 Runtime-generated (gitignored)
```
.devenv/
  worktrees/
    <env>/                 # git worktree
  instances/
    <env>/
      compose.yaml          # rendered compose
      instance.env          # resolved placeholder vars (optional)
      state.json            # instance metadata
      control.sock          # per-instance control socket (host)
```

---

## 4. Core CLI Commands

### 4.1 `devenv init`
**Purpose**: initialize a repo with `.devenv/` files.

**Flags**
- `--force`
- `--template <path>`

**Behavior**
- Create `.devenv/` if missing
- Create `.devenv/devenv-compose.yaml` template with a minimal agent host + example service
- Create `.devenv/config.yml` with minimal schema
- Add runtime dirs to `.gitignore` (or print instructions)
- Does not start anything

**Idempotency**: must not overwrite existing files unless `--force`.

---

### 4.2 `devenv up <env>`
**Purpose**: create or start an environment instance.

**Flags**
- `--from <ref>`
- `--agent "<prompt>"`
- `--agent-file <path>`
- `--no-agent`
- `--engine <docker|podman|nerdctl>`
- `--force`
- `--ports <placeholder=port,...>`

**Behavior**
1) Validate repo is initialized (`.devenv/config.yml` exists)
2) Normalize `<env>` (lowercase; spaces to `-`; collapse multiple `-`; trim `-`)
3) Validate normalized `<env>` matches `^[a-z0-9][a-z0-9-_]{0,62}$`
4) Worktree:
   - create if missing: `git worktree add .devenv/worktrees/<env> <ref>`
   - if exists: require clean unless `--force`
5) Port leasing:
   - for each placeholder in project config, lease a host port for this instance
   - allow explicit overrides via `--ports` if valid and available
   - stable across restarts (same instance gets same ports unless released)
6) Render compose template:
   - substitute placeholders with leased ports
   - write `.devenv/instances/<env>/compose.yaml`
   - optionally write `.devenv/instances/<env>/instance.env`
7) Start compose project:
   - unique project name derived from repo identity + env name
   - start stack in detached mode
8) Create per-instance control socket on host and mount into agent host container
9) If `--agent` supplied:
   - start agent inside agent host container with given prompt

**Output**
- status summary
- endpoints derived from placeholders (e.g., `frontend: http://localhost:<port>`)
- hint to attach agent if not started

---

### 4.3 `devenv down <env>`
**Purpose**: stop environment instance.

**Flags**
- `--purge`
- `--volumes`

**Behavior**
- Stop compose project for env
- If `--purge`:
  - remove containers/volumes (if `--volumes`)
  - delete `.devenv/instances/<env>/`
  - remove git worktree `.devenv/worktrees/<env>/`
  - release leased ports

---

### 4.4 `devenv list`
**Flags**
- `--json`

Lists known instances for this repo (based on `.devenv/instances/*/state.json`).

---

### 4.5 `devenv status <env>`
**Flags**
- `--json`
- `--verbose`

Shows:
- compose services and statuses (running/healthy/exited)
- agent status (running/last exit code/attached)
- port map

---

### 4.6 `devenv ports <env>`
**Flags**
- `--json`

Outputs:
- placeholder -> host port
- useful endpoints (optional mapping in config: “frontend = http://localhost:{FRONTEND_PORT}”)

---

## 5. Service Commands
All service operations are **two-argument**:
- `devenv service <op> <env> <svc> ...`

Where `<svc>` is a **service alias** from `.devenv/config.yml`, mapping to a compose service name.

### 5.1 `devenv service logs <env> <svc>`
**Flags**
- `--follow`
- `--tail N`
- `--since <duration>`
- `--json`

### 5.2 `devenv service restart <env> <svc>`

### 5.3 `devenv service exec <env> <svc> -- <cmd...>`
**Flags**
- `--workdir <path>`
- `--env KEY=VAL` (repeatable)
- `--user <uid:gid>`

Service groups (optional):
- `devenv service restart <env> app`
- `devenv service logs <env> app --follow`

---

## 6. Agent Commands
### 6.1 `devenv agent <env> "<prompt>"`
**Flags**
- `--agent <type>`
- `--provider <p>`
- `--model <m>`
- `--prompt-file <path>`
- `--prompt -`

Starts the agent inside the agent host service for that env and wires in internal MCP control.

Recommended: `devenv agent status <env>` / `devenv agent stop <env>`.

---

## 7. MCP Servers
### 7.1 Internal MCP (in-container): `devenv internal-mcp`
- stdio MCP server
- no `<env>` parameter
- tools:
  - `service.list` (list services)
  - `service.logs` (service logs)
  - `service.restart` (restart service)
  - `service.start` / `service.stop` (service lifecycle)
  - `ports.list` (leased ports)
  - `env.info` (env metadata)

### 7.2 Repo MCP (host): `devenv stdio`
**Flags**
- `--log-level`

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

## 8. Errors and Exit Codes
- `0` success
- `1` unexpected failure
- `2` invalid usage / bad flags / invalid env name
- `3` repo not initialized (`.devenv/config.yml` missing)
- `4` engine/compose backend failure
- `5` port lease failure / port conflicts
- `6` worktree conflict (dirty, exists, bad ref)
- `7` security violation (port exposure rule, scope)

Errors must include a short message and a next-step hint. If `<env>` is normalized, output both the input and normalized names.

---

## 9. Naming and Validation
- Normalize `<env>`: lowercase; spaces to `-`; collapse multiple `-`; trim leading/trailing `-`.
- Validate normalized `<env>` matches `^[a-z0-9][a-z0-9-_]{0,62}$`.
- Placeholders: `{[A-Z0-9_]+}` only.
- Service alias must exist in `.devenv/config.yml`.
- Unknown service aliases or placeholders must fail fast.

---

## 8. Safety & Scoping Requirements
- No container engine socket mounted in agent containers.
- Internal MCP scoped via per-instance socket + capability token.
- Strict port exposure enforcement when enabled.

---

## 9. Compatibility Notes
- Docker: baseline support
- Podman: supported via `podman compose` or docker-compat
- nerdctl: supported where `nerdctl compose` supports required features
