# Engineering Best Practices Audit — heimdall

| | |
|---|---|
| **Audit date** | 2026-11-19 |
| **Auditor** | Claude — gauge-repo skill |
| **Rubric version** | `item-credit-v1` — 2026-09-04 (`references/best-practices.md`) |
| **Audited ref** | branch `docs/ai-readiness` (PR head), on top of `main` @ `4f5e427` |

## Repo profile

Go service (`cmd/heimdall`, `cmd/deploydb`, 13 backend plugins built with `-buildmode=plugin`) plus a Next.js web UI in `web/`; the job-execution service the Pattern data platform submits work through. PostgreSQL is the only datastore (schema applied by `deploydb` from a list file). Published as the `patternoss/heimdall` image from `docker-image.yml`; where it is deployed is not recorded in this repo. About 20 contributors over the repo's life, 5–6 active. Owner verified with `gh api repos/patterninc/heimdall` → `patterninc/heimdall` (public), so inherited Pattern Wiz and Toolsmith controls apply (items 19, 20, 39, 47). GitHub rulesets `main + releases`, `require-pr-review` and `require-copilot-code-review` are active; none requires status checks and the repo has no `CODEOWNERS` although `main + releases` sets `require_code_owner_review`.

## Scorecard

| Metric | Value |
|--------|-------|
| **Critical gates** | **RED** |
| **Adjusted compliance** | **42.9%** |

Gates not Met: 16 (required CI: Gap), 24 (integration tests: Partial) and 40 (scoped secrets: Partial). Adjusted compliance is calculated independently:

`(16 Met + 0.5 × 10 Partial) / (49 total - 0 N/A) = 21.0 / 49 = 42.9%`

### Status totals

| Status | Items |
|--------|------:|
| Met | 16 |
| Partial | 10 |
| Gap | 23 |
| N/A | 0 |
| **Total** | **49** |

### Per-category breakdown

| Category | Met | Partial | Gap | N/A |
|----------|----:|--------:|----:|----:|
| Documentation & Context | 6 | 2 | 1 | 0 |
| Guardrails & Enforcement | 4 | 2 | 7 | 0 |
| Testing & Feedback Loops | 3 | 1 | 9 | 0 |
| Environment & Tooling | 3 | 4 | 6 | 0 |
| Agent dispatch | 0 | 1 | 0 | 0 |
| **Total** | **16** | **10** | **23** | **0** |

### Changes since last audit

`main` has no audit file (an earlier one was added in #144 and is no longer on `main`). Estimated state of `main` before this PR, from the same rubric: about 30.6% (11 Met, 8 Partial), gates RED. This is an estimate, not an audit. This PR is documentation only; the CI, tooling and test items above are recommendations, not changes in this PR.

## Documentation & Context

| # | Practice | Status | Evidence | Recommendation / rationale |
|---|----------|--------|----------|----------------------------|
| 1 | Skills / reusable prompt workflows | **Met** | `.agents/skills/heimdall-ui.md`, `heimdall-backend.md`, `heimdall-release.md` | — |
| 2 | AGENTS.md | **Met** | `AGENTS.md` (router, commands, conventions) | — |
| 3 | Architecture decision records | **Met** | `docs/adr/0001`–`0003` | Written retroactively and labelled as such. |
| 4 | Runbooks | **Met** | `docs/runbooks/` (deploy and rollback, stuck/AWOL jobs, cluster health, schema deploy) | `TODO(owner)` markers for deployment specifics. |
| 5 | API contract docs | **Partial** | `docs/openapi/heimdall.yaml` (all routes) | Bodies are generic objects and nothing checks the spec against the router. Add schemas and a parity test. |
| 6 | README with setup & run | **Met** | `README.md` | — |
| 7 | Changelog with migration notes | **Met** | `CHANGELOG.md` (v1.8.4–v1.10.0) | Earlier releases not backfilled. |
| 8 | On-call playbooks | **Partial** | `docs/oncall.md` | Triage table exists; paging, dashboards and escalation contacts are not in the repo. |
| 9 | CODEOWNERS | **Gap** | No `CODEOWNERS`; ruleset sets `require_code_owner_review` | Owner must name a team with write access. A team without write access would block all merges, so none was guessed. |

## Guardrails & Enforcement

| # | Practice | Status | Evidence | Recommendation / rationale |
|---|----------|--------|----------|----------------------------|
| 10 | Linters | **Partial** | `web/` ESLint (`web/eslint.config.mjs`); no Go linter | Add `golangci-lint` with the existing findings baselined. |
| 11 | Formatters | **Partial** | `web/` prettier scripts; no Go formatter check in CI | Add a `gofmt` check to CI (`starrocks/job_context.go` is not gofmt-clean today). |
| 12 | Type checking | **Met** | Go compiler; `web` script `type-check` | — |
| 13 | Pre-commit hooks | **Gap** | No `.pre-commit-config.yaml` | Add local hooks that mirror CI. |
| 14 | Commit message conventions | **Gap** | No commit or PR-title check | Add a Conventional Commits PR-title check. |
| 15 | Branch protection | **Met** | Ruleset `main + releases` (no deletion, no force push, PR review, squash only) | — |
| 16 | Required CI checks | **Gap** | No ruleset requires status checks | Require `build` (the `Build and Test` job) at minimum. Gate. |
| 17 | Dependency allow/deny lists | **Gap** | No deny-list or `go mod verify` in CI | Add a deny-list file and a check job. |
| 18 | License compliance scanning | **Gap** | No license scan | Add `go-licenses check` in CI. |
| 19 | Secret scanning | **Met** | Inherited Pattern Wiz policy | — |
| 20 | SAST gates | **Met** | Inherited Pattern Wiz policy | — |
| 21 | Max complexity limits | **Gap** | No complexity limit | Enable `gocyclo` (highest function today: 28) with the linter. |
| 22 | Import boundary enforcement | **Gap** | No import-boundary rule | Enable `depguard`; `pkg/plugin/runtime.go` importing `internal/pkg/aws` is an existing violation to baseline. |

## Testing & Feedback Loops

| # | Practice | Status | Evidence | Recommendation / rationale |
|---|----------|--------|----------|----------------------------|
| 23 | Unit tests | **Met** | 21 `_test.go` files; `build.yml` runs `./build.sh --test` | Coverage is low (7.9%); see 29. |
| 24 | Integration tests | **Partial** | `internal/pkg/janitor/stale_jobs_test.go` needs Postgres via `HEIMDALL_TEST_DATABASE`; `docs/runbooks/database-schema-deploy.md` explains how to run it | The test skips when the variable is unset and no CI job sets it, so it never runs in CI. Add a job with a Postgres service. Gate. |
| 25 | Snapshot / golden-file tests | **Gap** | No golden files | Add golden tests for result formatting and the Trino parser output. |
| 26 | Contract tests | **Gap** | No contract test | Compare the router with `docs/openapi/heimdall.yaml` in a test. |
| 27 | End-to-end tests | **Gap** | None for `web/` | Add Playwright against the smoke environment. |
| 28 | Visual regression tests | **Gap** | None | Add after Playwright. |
| 29 | Test coverage thresholds | **Gap** | No coverage threshold (measured 7.9% excluding generated grammar) | Enforce a floor just below the measured value. |
| 30 | Mutation testing | **Gap** | None | Add a scheduled informational workflow. |
| 31 | Load / performance benchmarks | **Gap** | No benchmarks | Add benchmarks for result decoding and run them on a schedule. |
| 32 | Flaky test quarantine | **Met** | `AGENTS.md` and `.agents/skills/heimdall-backend.md` document `t.Skip("flaky: <ticket>")` | — |
| 33 | Structured CI output | **Gap** | Test output is not uploaded as a machine-readable artifact | Run `go test -json` and upload it. |
| 34 | Deterministic test fixtures | **Met** | Janitor tests run in a rolled-back transaction with fixed rows; table-driven unit tests | — |
| 35 | Deploy smoke tests | **Gap** | No smoke test | Add a script that checks `/api/v1/clusters/health` and runs the ping job. |

## Environment & Tooling

| # | Practice | Status | Evidence | Recommendation / rationale |
|---|----------|--------|----------|----------------------------|
| 36 | Devcontainer | **Gap** | No devcontainer | Add `.devcontainer/devcontainer.json`. |
| 37 | One-command setup | **Gap** | No Makefile or setup script | Add one target that starts Postgres and applies the schema. |
| 38 | Seed scripts | **Partial** | `assets/databases/heimdall/data/*.sql` (status lookups), `configs/local.yaml` commands and clusters | No sample job history. |
| 39 | MCP servers | **Met** | Inherited Toolsmith-managed access | — |
| 40 | Scoped secrets per environment | **Partial** | `docker-image.yml` uses repo-level `DOCKERHUB_USERNAME` / `DOCKERHUB_TOKEN`; the only GitHub environment is `copilot` | Move the publish secrets into a protected environment and set `environment:` on the publish job. Needs an edit to an existing workflow. Gate. |
| 41 | Preview environments per PR | **Gap** | No per-PR environment | Start the server and Postgres per PR from CI. |
| 42 | Hot reload | **Gap** | No watch mode for the Go server | Add an `air` config. |
| 43 | Structured logging | **Gap** | `log.Println` plain text in `internal/pkg/rbac/ranger` and elsewhere | Move to a JSON logger with job and cluster IDs. Needs code changes. |
| 44 | Traces and metrics | **Partial** | Prometheus endpoint (`/metrics`), `go-telemetry` method metrics | No tracing evidence in this repo. |
| 45 | Feature flags | **Gap** | None | Needs code changes. |
| 46 | Database migration tooling | **Partial** | `cmd/deploydb` with ordered `heimdall.lst` | No versioned migrations or history table. |
| 47 | Dependency update automation | **Met** | Org-wide Wiz | — |
| 48 | Reproducible builds | **Met** | `go.sum`, `web/pnpm-lock.yaml`, `pnpm install --frozen-lockfile` in `build.sh` | — |

## Agent dispatch

| # | Practice | Status | Evidence | Recommendation / rationale |
|---|----------|--------|----------|----------------------------|
| 49 | Agent-dispatch manifest | **Partial** | `.agents/pattern-agents.json` has `schema_version`, `github.repo`, `skills` | Missing `clickup_list_id`, `slack_channel`, `datadog` (and `aws[]` if deployed to AWS). Owner must supply them. |

## Prioritized recommendations

1. **[S] Gap — required CI (item 16, gate):** a repo admin adds a ruleset requiring the `build` check (and each new check as it is added).
2. **[M] Partial — integration tests in CI (item 24, gate):** add a workflow with a Postgres service that applies the schema with `deploydb` and runs the janitor and database tests.
3. **[S] Partial — scoped secrets (item 40, gate):** put the DockerHub secrets used by `docker-image.yml` in a protected GitHub environment and set `environment:` on the publish job.
4. **[S] Gap — CODEOWNERS (item 9):** name the owning team (needs write access; the ruleset already sets `require_code_owner_review`).
5. **[S] Partial — agent manifest (item 49):** supply `clickup_list_id`, `slack_channel` and `datadog.service`.
6. **[S] Gap — linter, formatter, complexity, import rules (items 10, 11, 21, 22):** add `golangci-lint` and a `gofmt` check with today's findings baselined.
7. **[S] Gap — PR-title check, pre-commit, license scan, dependency deny-list (items 13, 14, 17, 18).**
8. **[M] Gap — coverage floor, benchmarks, CI test output (items 29, 31, 33).**
9. **[M] Gap — smoke test and per-PR environment (items 35, 41).**
10. **[M] Gap — contract test, golden files (items 25, 26); structured logging and feature flags (items 43, 45).**
11. **[L] Gap — Playwright e2e, visual regression, mutation testing (items 27, 28, 30).**

## Declined practices

None. No item is marked N/A.

## Beyond the checklist

- The janitor lease logic is table-tested against a real Postgres inside a rolled-back transaction (`internal/pkg/janitor/stale_jobs_test.go`).
- Per-area agent guidance already existed for the web app (`web/AGENTS.md`, `.agents/skills/heimdall-ui.md`).
