# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
# Build all binaries (requires UI to be built first)
make build

# Build individual components
make build-server   # builds UI first, then server binary -> dist/woodpecker-server
make build-agent    # -> dist/woodpecker-agent
make build-cli      # -> dist/woodpecker-cli

# UI only
make build-ui       # runs pnpm install + pnpm build in web/

# Dev UI (hot reload)
cd web && pnpm install --frozen-lockfile && pnpm dev
```

## Test Commands

```bash
# Run all Go tests
make test

# Run tests by component
make test-server
make test-agent
make test-cli
make test-lib
make test-server-datastore

# Run a single test
go test -run TestFunctionName -v ./server/pipeline/...

# UI tests
make test-ui
```

## Lint & Format

```bash
make lint           # golangci-lint (Go)
make lint-ui        # pnpm lint (Vue/TS)
make format         # gofumpt

# Install all dev tools
make install-tools
```

## Code Generation

```bash
make generate           # mockery + OpenAPI + go generate
make generate-openapi   # OpenAPI spec only
```

## Architecture

Woodpecker is a CI/CD system with three binaries sharing this monorepo:

- **server** (`cmd/server`) — HTTP API, webhook receiver, pipeline scheduler, gRPC server for agents
- **agent** (`cmd/agent`) — polls server via gRPC, executes pipeline steps via a backend
- **cli** (`cmd/cli`) — user-facing CLI client

### Request Flow

```
VCS Webhook → server/api → forge (webhook parsing)
                         → server/pipeline (compile + schedule)
                         → gRPC → agent/runner
                                  → pipeline/backend (Docker/K8s/local)
                                  → gRPC log stream back to server
```

### Key Packages

| Package | Role |
|---|---|
| `server/forge/` | VCS integrations — implements `forge.Forge` interface. One subdir per forge: `github`, `gitea`, `forgejo`, `gitlab`, `bitbucket`, `bitbucketdatacenter` |
| `server/model/` | All data models (Repo, Pipeline, Workflow, Step, User, Secret, etc.) |
| `server/store/` | DB abstraction layer; `datastore/` has XORM-based SQLite/PostgreSQL/MySQL implementations |
| `server/pipeline/` | Pipeline lifecycle: compile, schedule, process events, step updates |
| `server/api/` | REST API handlers (Gin framework) |
| `server/router/` | HTTP routing setup |
| `server/rpc/` | gRPC server-side implementation (agent ↔ server protocol) |
| `pipeline/frontend/` | `.woodpecker.yml` parsing and compilation into internal step graph |
| `pipeline/backend/` | Execution backends — implements `backend.Backend` interface: `docker`, `kubernetes`, `local` |
| `agent/runner/` | Agent-side pipeline execution loop |
| `web/src/` | Vue 3 + TypeScript frontend (Vite, Pinia, WindiCSS) |

### Extension Points

- **New VCS forge**: implement `server/forge/forge.go` `Forge` interface, add to `server/forge/setup/`
- **New execution backend**: implement `pipeline/backend/backend.go` `Backend` interface
- **New API endpoints**: add handler in `server/api/`, register route in `server/router/`
- **DB schema change**: add migration in `server/store/datastore/migration/`
- **Pipeline syntax**: extend parsing in `pipeline/frontend/yaml/`

### Module Path

`go.woodpecker-ci.org/woodpecker/v3` — use this in all import paths.

### Go Version

Go 1.25+ required (see `go.mod`). Server uses `CGO_ENABLED=1` (for SQLite); agent/CLI use `CGO_ENABLED=0`.

### Frontend

UI lives in `web/`. Uses pnpm (not npm/yarn). Package manager: pnpm with frozen lockfile in CI.

### gRPC Proto

Proto definitions at `pipeline/rpc/proto/`. After editing, run `make` to regenerate. Requires `protoc-gen-go` and `protoc-gen-go-grpc` (installed via `make install-tools`).
