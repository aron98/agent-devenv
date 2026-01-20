# AGENTS

## Purpose
- Run coding agents inside per-environment containers.
- Keep agent control scoped to a single environment.
- Avoid host engine socket access in containers.

## How agents run
- The agent host service is configured in `.devenv/config.yml`.
- Start an agent with `devenv up <env> --agent "<prompt>"` or `devenv agent <env> "<prompt>"`.
- Prompts are recorded at `.devenv/instances/<env>/prompt.md`.

## MCP and control scope
- In-container tools use `devenv internal-mcp` and are scoped to the current env.
- Host orchestration uses `devenv stdio` scoped to the repo.
- Internal MCP can list services, read logs, and restart services only for its env.

## Safety and constraints
- Never mount the container engine socket into agent containers.
- Use per-instance control sockets and capability tokens.
- Enforce port exposure rules based on placeholders in `.devenv/config.yml`.
- Agents must not control services outside their environment.

## Operational notes
- Each env has a git worktree under `.devenv/worktrees/<env>`.
- Rendered Compose files live under `.devenv/instances/<env>/compose.yaml`.
- Use `devenv down <env> [--purge]` to stop or remove an environment.
