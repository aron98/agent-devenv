# Testing

## Automated tests
- `go test ./...`

## Smoke test (manual)
1. `go run ./cmd/devenv help`
   - Expect usage output on stdout and exit code 0.
2. `go run ./cmd/devenv`
   - Expect a "no command provided" error and exit code 2 (`echo $?`).
3. `go run ./cmd/devenv init`
   - Expect a "not implemented" error and exit code 1 (`echo $?`).
