# Engineering Best Practices Audit — heimdall

| | |
|---|---|
| **Audit date** | 2026-09-09 |
| **Auditor** | Claude — gauge-repo skill |
| **Rubric version** | `item-credit-v1` — 2026-09-04 (`references/best-practices.md`) |

## Repo profile

Heimdall is a data orchestration and job execution platform: a Go 1.25 backend that exposes a REST API and a Next.js 16 / React 19 web UI (`web/`, `@patterninc/react-ui`, pnpm 10.28.2). It uses PostgreSQL 16 for persistence (DDL under `assets/databases/heimdall/` applied by the in-repo `deploydb` binary via `heimdall.lst`), Prometheus for metrics, and a plugin-based executor for Snowflake, Spark, EKS, EMR, ECS Fargate, Glue, DynamoDB, Trino, ClickHouse, StarRocks, Postgres, and Shell. It has a clear AWS footprint (S3, ECS, EKS, EMR, Glue, DynamoDB, STS via aws-sdk-go-v2). The image `patternoss/heimdall` is built multi-arch and pushed to DockerHub from `main`/`v*` tags; the Backstage descriptor labels it `Owner: dev-data-acquisition`, `System: data`, `Environment: stage`. The repo has ~15–20 unique contributors and active development. GitHub owner is `patterninc` (verified via `gh repo view`), so inherited Pattern Wiz (items 19, 20, 47) and Toolsmith (item 39) controls apply.

## Scorecard

| Metric | Value |
|--------|-------|
| **Critical gates** | **RED** |
| **Adjusted compliance** | **40.8%** |

Critical gates are RED because required CI (16), integration tests (24), and scoped-per-environment secrets (40) are Partial. Adjusted compliance is calculated independently:

`(15 Met + 0.5 × 10 Partial) / (49 total - 0 N/A) = 40.8%`

### Status totals

| Status | Items |
|--------|------:|
| Met | 15 |
| Partial | 10 |
| Gap | 24 |
| N/A | 0 |
| **Total** | **49** |

### Per-category breakdown

| Category | Met | Partial | Gap | N/A |
|----------|----:|--------:|----:|----:|
| Documentation & Context | 3 | 1 | 5 | 0 |
| Guardrails & Enforcement | 4 | 4 | 5 | 0 |
| Testing & Feedback Loops | 1 | 3 | 9 | 0 |
| Environment & Tooling | 7 | 2 | 4 | 0 |
| Agent dispatch | 0 | 0 | 1 | 0 |
| **Total** | **15** | **10** | **24** | **0** |

## Documentation & Context

| # | Practice | Status | Evidence | Recommendation / rationale |
|---|----------|--------|----------|----------------------------|
| 1 | Skills / reusable prompt workflows | **Met** | `.agents/skills/heimdall-ui.md`, `.agents/personas/frontend-reviewer.md`, `AGENTS.md` routes to them | — |
| 2 | AGENTS.md | **Met** | `AGENTS.md` — router pointing to per-task skills and personas | — |
| 3 | Architecture decision records | **Gap** | No `docs/adr/` or equivalent directory | Record load-bearing decisions (plugin loading model, sync vs. async execution, cluster/command matching, `deploydb` vs. an off-the-shelf migrator) as dated ADRs under `docs/adr/`. |
| 4 | Runbooks | **Gap** | No `docs/runbooks/`; only inline README instructions | Add operator runbooks for common tasks (rotate DockerHub creds, wipe stuck async jobs, restore Postgres, roll back a bad plugin release). |
| 5 | API contract docs (OpenAPI / protobuf) | **Partial** | README documents endpoints in a table but no machine-readable spec | Publish an OpenAPI 3 spec generated from or checked against the Gorilla router in `internal/pkg/heimdall/heimdall.go`; wire client generation for the web UI. |
| 6 | README with setup and run instructions | **Met** | `README.md` covers clone, `docker compose up --build -d`, curl example, and `build.sh` flag matrix | — |
| 7 | Changelog with migration notes | **Gap** | No `CHANGELOG.md`; GitHub tags `v*` exist but no per-release notes with upgrade steps | Adopt Keep-a-Changelog (or generate from Conventional Commits) and call out breaking schema/config changes for consumers of `deploydb`. |
| 8 | On-call playbooks | **Gap** | No `docs/oncall/` or equivalent | Add a playbook covering "jobs stuck", "cluster health probes failing", "plugin .so load failure", and "Postgres unreachable". |
| 9 | CODEOWNERS | **Gap** | `.github/CODEOWNERS` absent; ruleset `main + releases` has `require_code_owner_review: true` but there is no CODEOWNERS file for it to consult | Add `CODEOWNERS` mapping `web/` to the frontend owners, `plugins/*/` to plugin maintainers, and `assets/databases/**` to data platform. |

## Guardrails & Enforcement

| # | Practice | Status | Evidence | Recommendation / rationale |
|---|----------|--------|----------|----------------------------|
| 10 | Linters | **Partial** | `web/eslint.config.mjs` (Next.js core-web-vitals + TS presets); no `.golangci.yml`; CI runs only `go test` + `next build` | Add `.golangci.yml` and run `golangci-lint run` in `.github/workflows/build.yml`; run `pnpm lint` for `web/` too. |
| 11 | Formatters | **Partial** | Prettier configured in `web/package.json` (`format`, `prettier-check`); no `gofmt`/`goimports` enforcement in CI | Wire `pnpm prettier-check` and `gofmt -l` (fail on diff) into the Build and Test workflow. |
| 12 | Type checking | **Met** | `web/tsconfig.json` sets `strict: true`; Go is statically typed and enforced by `go build` in CI | — |
| 13 | Pre-commit hooks | **Gap** | No `.pre-commit-config.yaml`, `.githooks/`, or husky install | Add a lightweight pre-commit that runs `gofmt`, `golangci-lint --fast`, and `pnpm prettier-check` on staged files. |
| 14 | Commit message conventions | **Partial** | Recent history mixes Conventional Commits with ad-hoc casing (`FEAT:`, `feat:`, `fix:`, `CHORE:`, plain "Starrocks-db-connection-close"); no commitlint | Standardize on lowercase Conventional Commits and enforce with commitlint (or a GitHub Action) on PR titles. |
| 15 | Branch protection | **Met** | Rulesets `require-pr-review` and `main + releases` block deletion and non-fast-forward, require PR with 1 approval and stale-review dismissal | — |
| 16 | Required CI checks before merge | **Partial** | `.github/workflows/build.yml` runs `build.sh --go --ui --test` on `pull_request`, but neither ruleset lists a `required_status_checks` rule | Add a `required_status_checks` rule referencing the Build and Test job so merges block on red CI. |
| 17 | Dependency allow-lists / deny-lists | **Gap** | No policy file; `go.mod` and `web/package.json` are unrestricted | Introduce a lightweight policy (e.g. `go mod tidy` verification job + a curated deny-list for known-bad JS packages) or adopt an internal policy doc referenced from AGENTS.md. |
| 18 | License compliance scanning | **Gap** | No `licensei`, `license-checker`, or FOSSA job | Add a license scan (e.g. `go-licenses check ./...` and `pnpm license-checker`) to the PR workflow. |
| 19 | Secret scanning | **Met** | Inherited Pattern Wiz policy (verified `patterninc` origin) | — |
| 20 | SAST / static analysis gates | **Met** | Inherited Pattern Wiz policy (verified `patterninc` origin) | — |
| 21 | Max complexity limits | **Gap** | No cyclomatic complexity or function-length rules configured | Enable `gocyclo`/`funlen` in `.golangci.yml` when it lands (see item 10). |
| 22 | Import boundary enforcement | **Gap** | `pkg/`, `internal/`, `plugins/` split by convention only | Add `depguard` or `go-arch-lint` rules to prevent `plugins/*` from importing `internal/*` and to keep `pkg/*` free of `internal/*`. |

## Testing & Feedback Loops

| # | Practice | Status | Evidence | Recommendation / rationale |
|---|----------|--------|----------|----------------------------|
| 23 | Unit tests | **Met** | `pkg/result/result_test.go`, `pkg/result/column/type_test.go`, `internal/pkg/sql/parser/trino/tests/*`, `internal/pkg/rbac/ranger/tests/*`, `internal/pkg/object/command/sparkeks/*_test.go`, `internal/pkg/object/command/ecs/ecs_test.go` | — |
| 24 | Integration tests | **Partial** | `docker-compose.yaml` provisions Postgres; command-object tests exercise per-plugin logic — but there is no test target that starts the API against a live DB and exercises `/api/v1/job` end-to-end | Add an `integration` build tag / `make integration` target that boots `docker compose`, runs `deploydb`, and drives the REST API through a happy-path job lifecycle. |
| 25 | Snapshot / golden-file tests | **Gap** | Parser tests hard-code expected structs; no `testdata/` golden fixtures | Convert the Trino parser expected outputs to golden JSON under `testdata/` so diffs are easy to review. |
| 26 | Contract tests (Pact) | **Gap** | REST API has at least two consumers (the bundled web UI and external `curl` clients) with no shared contract | Generate an OpenAPI spec (item 5), then have `web/` regenerate its client from it in CI so drift fails the build; consider a Pact broker if third parties integrate. |
| 27 | End-to-end tests | **Gap** | Web UI exists at `web/`; no Playwright/Cypress harness | Add a small Playwright suite that logs in as a test user and walks the Jobs → Job Detail → Cancel flow against `docker compose`. |
| 28 | Visual regression tests | **Gap** | UI ships real user-facing components (`@patterninc/react-ui`) with no screenshot diffing | If the Playwright suite (item 27) lands, layer Playwright screenshot comparisons on the key pages. |
| 29 | Test coverage thresholds | **Gap** | `build.sh --test` runs `go test` without `-cover`/`-coverprofile`; no CI threshold | Emit `go test -coverprofile` in CI and fail the job below a minimum (start at existing baseline). |
| 30 | Mutation testing | **Gap** | Not configured | Trial `gremlins` on `pkg/result` and `internal/pkg/sql/parser/trino` where tests are densest. |
| 31 | Load / performance benchmarks | **Gap** | No `*_bench_test.go` files; no k6/vegeta harness | Add Go `Benchmark*` for the plugin dispatch hot path and a k6 script that submits N sync + N async jobs. |
| 32 | Flaky test quarantine | **Gap** | No quarantine or retry mechanism | Once integration tests (item 24) land, add a documented convention for `t.Skip` + linked ticket for known-flaky cases. |
| 33 | Structured CI output | **Partial** | `go test` and `next build` write plain text to Actions logs; no JUnit XML uploaded | Emit JUnit via `gotestsum --junitfile` and upload as an artifact so PR annotations attach to failing tests. |
| 34 | Deterministic test fixtures | **Partial** | Parser tests use fixed input strings; `internal/pkg/object/command/sparkeks/entrypoint_test.go` uses inline data; no shared fixture layer with clocks/UUID seeding | Introduce a `testfixtures` package that centralizes stubbed `time.Now`, deterministic UUIDs, and canned Postgres seed data. |
| 35 | Smoke tests for deploys | **Gap** | `.github/workflows/docker-image.yml` builds and pushes the multi-arch image but does not `docker run` it or hit `/api/v1/clusters/health` post-push | Add a final job that pulls the freshly pushed tag, runs it against a throwaway Postgres, and probes `/api/v1/clusters/health` + a ping job. |

## Environment & Tooling

| # | Practice | Status | Evidence | Recommendation / rationale |
|---|----------|--------|----------|----------------------------|
| 36 | Devcontainer config | **Gap** | No `.devcontainer/` | Add a minimal devcontainer pinned to `golang:1.25.0` and `node:20-bookworm` (matching `Dockerfile`) with pnpm and Docker-outside-of-Docker for `docker compose`. |
| 37 | One-command setup | **Met** | `docker compose up --build -d` starts Heimdall + Postgres; `./build.sh --go --ui --test` covers the full build/test loop | — |
| 38 | Seed scripts for local databases | **Met** | `assets/databases/heimdall/data/{command,cluster,job}_statuses.sql` applied by `cmd/deploydb` from `heimdall.lst` | — |
| 39 | MCP servers for external tools | **Met** | Toolsmith-managed MCP access (inherited Pattern control); `.agents/skills/heimdall-ui.md` documents the chrome-devtools MCP hook | — |
| 40 | Scoped secrets per environment | **Partial** | `docker-compose.yaml` reads `AWS_*` from the host shell; `.github/workflows/docker-image.yml` uses `DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN` as flat repo secrets — no GitHub Environments split for stage vs. prod, no separation of scanning vs. push credentials | Move DockerHub push into a `production` GitHub Environment with its own reviewers; document expected env-var scoping for `AWS_*` in the README. |
| 41 | Preview environments per PR | **Gap** | No preview deployment workflow | Optional given internal ops-tool posture; if adopted, an ECS Fargate task or Fly.io per-PR deploy off the merged image would be enough. |
| 42 | Hot-reload / watch mode | **Met** | `web/package.json` `dev` script (`PORT=4000 next dev`); Go binaries rebuild via `go run` / `./build.sh --go` | — |
| 43 | Structured logging (JSON) | **Gap** | Async worker uses `fmt.Println(...)` (`internal/pkg/heimdall/jobs_async.go:164,171,176,181`); no `log/slog`, `zap`, or `zerolog` handler | Standardize on `log/slog` with a JSON handler, thread `job_id`/`cluster_id`/`user` as attributes, and remove ad-hoc `fmt.Println` calls. |
| 44 | Observable traces and metrics | **Partial** | Prometheus metrics endpoint exposed via reverse proxy at `/metrics` (`internal/pkg/heimdall/metrics.go`, `PROMETHEUS_ADDRESS` in `docker-compose.yaml`); no OpenTelemetry traces, no application-level counters or histograms defined in-repo | Instrument the plugin dispatch and job lifecycle with OpenTelemetry (traces + metrics) and export via OTLP; keep the Prometheus scrape as a supplemental sink. |
| 45 | Feature flags with local overrides | **Gap** | No feature-flag SDK; toggles happen via YAML config edits | Introduce a lightweight flag layer (Unleash SDK or config-driven bool map with hot reload) for risky plugin rollouts. |
| 46 | Database migration tooling | **Met** | `cmd/deploydb/deploydb.go` applies ordered SQL files from `assets/databases/heimdall/build/heimdall.lst`; individual DDL under `assets/databases/heimdall/tables/*.sql` | — |
| 47 | Dependency update automation | **Met** | Inherited Pattern Wiz policy (verified `patterninc` origin) | — |
| 48 | Reproducible builds (lockfiles) | **Met** | `go.sum` committed; `web/pnpm-lock.yaml` committed; `Dockerfile` pins `golang:1.25.0` and `node:20-bookworm`; `build.sh` calls `pnpm install --frozen-lockfile` | — |

## Documentation & Context (agent dispatch)

| # | Practice | Status | Evidence | Recommendation / rationale |
|---|----------|--------|----------|----------------------------|
| 49 | Agent-dispatch manifest | **Gap** | `.agents/` exists but contains only `skills/` and `personas/` — no `pattern-agents.json`/`.yml`/`.yaml` with `schema_version`, `github.repo`, `clickup_list_id`, `slack_channel`, `datadog.service`, `skills.plugins`, or the required `aws[]` array (repo deploys to AWS) | Add `.agents/pattern-agents.json` covering GitHub (`patterninc/heimdall`), the data-acquisition ClickUp list, the team Slack channel, the Datadog service name, active skills, and one `aws[]` entry per deployed account (mark the primary as `default: true`). |

## Prioritized recommendations

1. **[S] Partial — required CI (item 16, RED gate):** Add a `required_status_checks` rule to the `main + releases` ruleset that lists the Build and Test job so merges block on red CI.
2. **[M] Partial — integration tests (item 24, RED gate):** Add an `integration` build tag / target that boots `docker compose`, runs `deploydb`, and drives `/api/v1/job` end-to-end.
3. **[S] Partial — scoped secrets per environment (item 40, RED gate):** Move DockerHub push into a `production` GitHub Environment with its own reviewers and document `AWS_*` env-var scoping.
4. **[S] Gap — CODEOWNERS (item 9):** Add `.github/CODEOWNERS` so the existing `require_code_owner_review: true` rule can actually enforce ownership.
5. **[S] Gap — agent-dispatch manifest (item 49):** Create `.agents/pattern-agents.json` with GitHub, ClickUp, Slack, Datadog, skills, and per-account `aws[]` metadata.
6. **[S] Gap — changelog (item 7):** Adopt Keep-a-Changelog and start populating it from the existing `v*` tags.
7. **[S] Gap — on-call playbooks (item 8):** Add `docs/runbooks/` (or `docs/oncall/`) covering stuck-jobs, plugin `.so` load failures, DockerHub cred rotation, and Postgres recovery.
8. **[S] Gap — ADRs (item 3):** Capture the plugin-loading model, sync/async execution split, and `deploydb`-vs-off-the-shelf-migrator decisions under `docs/adr/`.
9. **[S] Partial — Go linter (item 10):** Add `.golangci.yml` and run `golangci-lint run` in `.github/workflows/build.yml`.
10. **[S] Partial — formatter enforcement (item 11):** Wire `gofmt -l` (fail on diff) and `pnpm prettier-check` into the Build and Test workflow.
11. **[S] Partial — commit conventions (item 14):** Standardize on lowercase Conventional Commits and enforce with commitlint on PR titles.
12. **[S] Gap — structured logging (item 43):** Migrate `fmt.Println` in `internal/pkg/heimdall/jobs_async.go` to `log/slog` with a JSON handler and structured attributes.
13. **[M] Partial — structured CI output (item 33):** Emit JUnit via `gotestsum --junitfile` and upload as a workflow artifact.
14. **[S] Gap — smoke tests for deploys (item 35):** Add a job that pulls the freshly pushed image and probes `/api/v1/clusters/health` and a ping job.
15. **[S] Gap — coverage thresholds (item 29):** Emit `go test -coverprofile` and fail below an established baseline.
16. **[M] Gap — E2E tests (item 27):** Add a small Playwright suite covering the Jobs list → detail → cancel flow.
17. **[M] Partial — OpenAPI (item 5) + contract tests (item 26):** Publish an OpenAPI 3 spec and have `web/` regenerate its client from it in CI.
18. **[M] Gap — pre-commit hooks (item 13):** Add a `.pre-commit-config.yaml` running `gofmt`, `golangci-lint --fast`, and `pnpm prettier-check`.
19. **[M] Gap — license compliance (item 18):** Add `go-licenses check ./...` and `pnpm license-checker` gates.
20. **[M] Gap — devcontainer (item 36):** Add a `.devcontainer/` pinned to the `Dockerfile`'s Go and Node versions.
21. **[M] Partial — deterministic fixtures (item 34):** Introduce a `testfixtures` package centralizing clocks, UUID seeds, and Postgres seed data.
22. **[M] Partial — observability (item 44):** Instrument the plugin dispatch and job lifecycle with OpenTelemetry traces and application-level metrics.
23. **[M] Gap — import boundary enforcement (item 22):** Add `depguard` rules to prevent `plugins/*` from reaching into `internal/*`.
24. **[L] Gap — visual regression (item 28):** Once Playwright lands, add screenshot diffing on the primary pages.

## Declined practices

No items were marked Not applicable. Every checklist item is at least in scope for this profile: the repo is a deployed service with a REST API, a bundled web UI, a database, an AWS footprint, and multiple contributors, so items that a narrower profile could decline (contract tests, E2E, visual regression, preview envs, playbooks, CODEOWNERS, migration tooling, structured logging) are all applicable here.

## Beyond the checklist

- **Router-style AGENTS.md** delegates task-specific guidance to `.agents/skills/` and `.agents/personas/` so the top-level file stays skimmable while the deep context is loaded on demand.
- **Plugin architecture with per-plugin READMEs** — the main README links out to each plugin's own README (`plugins/{ping,shell,glue,dynamo,snowflake,spark,sparkeks,trino,clickhouse,ecs,postgres}/README.md`), keeping surface-level docs close to the code they describe.
- **Config-driven job attributes** (documented in `README.md`) let plugins surface links and metadata to the UI via Go `text/template` without UI code changes — a genuine agent-friendliness win.
- **Backstage descriptor** (`backstage.yaml`) tags ownership (`dev-data-acquisition`), system (`data`), cost center, and environment, so service catalogs pick up the repo automatically.
- **Multi-arch Docker builds on native runners** — `.github/workflows/docker-image.yml` builds `linux/amd64` on `ubuntu-latest` and `linux/arm64` on `ubuntu-24.04-arm` (no QEMU) and stitches a manifest list, which is faster and more reliable than the common single-runner + QEMU pattern.
