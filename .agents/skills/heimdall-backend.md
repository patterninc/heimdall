# Skill: Heimdall Go backend changes

Agent-agnostic guide for changes under `cmd/`, `internal/`, `pkg/` and `plugins/`. UI-only changes use [`heimdall-ui.md`](heimdall-ui.md).

## Loop
1. Read the area first: server and routes in `internal/pkg/heimdall/heimdall.go`; job lifecycle in `internal/pkg/heimdall/job*.go`; janitor in `internal/pkg/janitor`; backends in `plugins/<name>` and `internal/pkg/object/command/<name>`.
2. Make the change. Keep `pkg/` free of imports from `internal/` and `plugins/` (a convention; `pkg/plugin/runtime.go` is a known exception).
3. Gate locally:
   - `gofmt -l cmd internal pkg plugins` (ignore the generated `internal/pkg/sql/parser/trino/grammar/`)
   - `./build.sh --test` (what the `Build and Test` CI job runs)
   - when the change touches SQL, `internal/pkg/janitor` or `internal/pkg/database`: start Postgres (`docker compose up -d postgres`, localhost:5433), apply the schema as in [`docs/runbooks/database-schema-deploy.md`](../../docs/runbooks/database-schema-deploy.md), and run `HEIMDALL_TEST_DATABASE=<connection string> go test ./internal/pkg/janitor/... ./internal/pkg/database/...`. Those tests skip silently when the variable is unset.
4. A plugin or startup change is not verified by unit tests alone: run the server against Postgres and submit the ping job from the README. Go plugins do not load on macOS; use Linux or a container.

## Rules that are easy to miss
- **New or changed route:** update `docs/openapi/heimdall.yaml` in the same change.
- **Schema change:** add or edit a SQL file under `assets/databases/heimdall/` and list it in `build/heimdall.lst`; it must be re-runnable.
- **Release-visible change:** add an entry under `[Unreleased]` in `CHANGELOG.md`, with a **Migration** note if schema, config or API changed.
- **Plugins** build with `-buildmode=plugin` and must use the same toolchain and dependency versions as the server; `go build ./...` printing `function main is undeclared` for `plugins/*` is expected. Use `./build.sh --go`.
- **Secrets:** never put credentials in config files, tests or logs; API responses redact `password`, `private_key`, `token` in `context`.
- **Decisions:** record structural decisions as a new ADR in `docs/adr/`.

## Flaky tests
Do not delete or silently skip a flaky test. Mark it with `t.Skip("flaky: <ticket>")` so it is searchable, and link the ticket in the PR.
