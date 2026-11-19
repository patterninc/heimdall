# AGENTS.md

Agent-agnostic guidance for working in this repo. Human docs live in `README.md`.

This file is a **router** — deeper agent context lives in the **`.agents/`** folder, loaded on
demand. Read the matching entry before starting a task it covers.

## Task guides — `.agents/skills/`

| When the task is…                                                                                                                                                   | Read                                                              |
| -------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------- |
| **UI-only** change to the web app (`web/`) — buttons/links, labels, layout, responsiveness, tooltips, table columns, rendering data the API already returns | [`.agents/skills/heimdall-ui.md`](.agents/skills/heimdall-ui.md) |
| **Any** change to the web app (`web/`) — points at the Next.js docs matching the installed version                                                                  | [`web/AGENTS.md`](web/AGENTS.md)                                  |

## Backend (Go) — read before any non-UI change

**Layout.** Binaries in `cmd/` (`heimdall` server, `deploydb` schema deployer, `result` CLI);
server in `internal/pkg/heimdall/` (routes in `heimdall.go`, SQL in `queries/`); public types in
`pkg/`; schema in `assets/databases/heimdall/`.

**Build and test.** `./build.sh --go --ui --test` is what CI runs, and it must pass to merge. Flags
are independent; `--test` fails on unformatted Go (`gofmt -l`), then runs
`go test ./internal/pkg/... ./pkg/...`.

**Integration tests** skip locally unless their variable is set; CI sets both. After
`docker compose up --build -d`:

- `HEIMDALL_TEST_DATABASE=postgres://heimdall:heimdall@localhost:5433/heimdall?sslmode=disable`
  runs the Postgres-backed tests (they roll back what they write).
- `HEIMDALL_TEST_URL=http://127.0.0.1:9090` runs `internal/pkg/apitest/`, which submits a ping job
  to the running server and reads its result. Add API-level cases there.

**Command plugins.** A thin `plugins/<name>/<name>.go` (`package main`, exports `New`) wraps the
implementation in `internal/pkg/object/command/<name>/`, which implements `plugin.Handler` from
`pkg/plugin` and optionally `HealthChecker`, `ResultStreamer` or `Resumer`. Document each in
`plugins/<name>/README.md`.

**Schema.** `deploydb` re-applies every file listed in `assets/databases/heimdall/build/heimdall.lst`
on each start, so DDL must be idempotent (`if not exists`, guarded `do $$` blocks). List new files
there.

**Rules.**

- Never edit generated code: `internal/pkg/sql/parser/trino/grammar/` (ANTLR, from the `.g4`
  files) and `internal/pkg/rbac/ranger/mocks/` (mockery). Both are marked in `.gitattributes`.
- Run `gofmt -w` on every Go file you touch; the build rejects unformatted code.
- Context keys `password`, `private_key` and `token` are redacted in API responses; add any other
  secret-bearing key to `sensitive_context_keys`.
- PR titles use Conventional Commits; merges are squash-only.

## Personas — `.agents/personas/`

Adopt a persona when it fits the task:

| Persona           | Use for                 | Profile                                                                           |
| ----------------- | ----------------------- | --------------------------------------------------------------------------------- |
| Frontend reviewer | UI-only `web/` changes | [`.agents/personas/frontend-reviewer.md`](.agents/personas/frontend-reviewer.md) |

## Setup (one-time, per machine)

Some guides need a tool your agent provides — e.g. UI browser verification uses the
**chrome-devtools MCP**. Two ways it gets set up:

- **On demand:** the guide's **Step 0** fires when a task first needs the tool.
- **Proactively:** ask your agent to **"set up the &lt;name&gt; skill"** (e.g. *"set up the
  heimdall-ui skill"*).
