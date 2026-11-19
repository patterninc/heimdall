# AGENTS.md

Agent-agnostic guidance for working in this repo. Human docs live in `README.md`.

This file is a **router** — deeper agent context lives in the **`.agents/`** folder, loaded on
demand. Read the matching entry before starting a task it covers.

## Task guides — `.agents/skills/`

| When the task is…                                                                                                                                                   | Read                                                              |
| -------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------- |
| **UI-only** change to the web app (`web/`) — buttons/links, labels, layout, responsiveness, tooltips, table columns, rendering data the API already returns | [`.agents/skills/heimdall-ui.md`](.agents/skills/heimdall-ui.md) |
| **Any** change to the web app (`web/`) — points at the Next.js docs matching the installed version                                                                  | [`web/AGENTS.md`](web/AGENTS.md)                                  |

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

## Backend, release and ops guides

| When the task is…                                              | Read                                                                       |
| -------------------------------------------------------------- | -------------------------------------------------------------------------- |
| A Go change under `cmd/`, `internal/`, `pkg/` or `plugins/`    | [`.agents/skills/heimdall-backend.md`](.agents/skills/heimdall-backend.md) |
| Cutting or rolling back a release                              | [`.agents/skills/heimdall-release.md`](.agents/skills/heimdall-release.md) |
| Operating Heimdall (stuck jobs, cluster health, schema deploy) | [`docs/oncall.md`](docs/oncall.md), [`docs/runbooks/`](docs/runbooks/README.md) |
| Why something is built the way it is                           | [`docs/adr/`](docs/adr/README.md)                                          |
| The HTTP API                                                   | [`docs/openapi/heimdall.yaml`](docs/openapi/heimdall.yaml)                 |

## Commands

- `./build.sh --test` — Go unit tests (what the `Build and Test` CI job runs). `./build.sh --go` builds binaries and plugins; `./build.sh --ui` builds the web app.
- `HEIMDALL_TEST_DATABASE=<connection string> go test ./internal/pkg/janitor/... ./internal/pkg/database/...` — database-backed tests; they skip when the variable is unset. Schema setup: [`docs/runbooks/database-schema-deploy.md`](docs/runbooks/database-schema-deploy.md).
- `docker compose up postgres` — local Postgres on `localhost:5433`.

## Conventions

- A route change updates `docs/openapi/heimdall.yaml`. A release-visible change updates `CHANGELOG.md`.
- A flaky test is marked `t.Skip("flaky: <ticket>")`, never deleted.
- `go build ./...` reports `function main is undeclared` for `plugins/*`; that is expected, build with `./build.sh --go`.
