# 0001. Backends are Go plugins loaded at runtime

- Status: Accepted (retroactive, 2026-11-19)
- Evidence: `build.sh` (`go build -buildmode=plugin ... plugins/<name>/<name>.go`), `plugins/*`, `configs/local.yaml` (`plugin_directory`, `plugin: ping`), `pkg/plugin`

## Context
Heimdall submits work to many backends (Trino, Snowflake, Spark, ECS, Glue, ...). Each needs its own client libraries and credentials handling.

## Decision
Each backend is a separate Go package under `plugins/` built as a `.so` into `dist/plugins/` and loaded from `plugin_directory` at start. A command in config
selects a plugin by name (`plugin: ping`). Shared contracts live in `pkg/plugin`.

## Consequences
- Adding a backend does not touch the server code path. (inferred rationale)
- Plugins must be built with the same Go toolchain and dependency versions as the server, so `build.sh` builds both together and `go build ./...` reports `function main is undeclared` for plugin packages. That message is expected; use `./build.sh --go`.
- Plugin packages are loaded into the server process, so a panic in a plugin affects the whole agent.
