# Changelog

Format: [Keep a Changelog](https://keepachangelog.com/). Versions are git tags `vX.Y.Z`. Every release that changes the database schema, the config file
or the API contract carries a **Migration** note. Entries before v1.8.4 are not backfilled.

## [Unreleased]
- Docs: runbooks, ADRs, on-call playbook, OpenAPI spec (`docs/openapi/heimdall.yaml`), agent skills. No runtime change.

## [1.10.0] - 2026-11-16
- Resumable jobs (sparkeks, ecs) are released after `stale_job` (45 s) and failed as AWOL only after `awol_timeout` (#164).
- Resume sparkeks and ecs jobs on another agent instead of failing them as AWOL (#162).
- Graceful shutdown drains in-flight jobs on SIGTERM, up to `HEIMDALL_SHUTDOWN_TIMEOUT` seconds (default 30) (#161).
- Web UI upgraded to Node 26, TypeScript 7, Next 16.4, pnpm 12 (#159). Image builds on Go 1.27.3 (#157).
- **Migration:** none to the schema or config file. Set the container stop timeout at or above `HEIMDALL_SHUTDOWN_TIMEOUT` so the drain can finish. See [ADR 0003](docs/adr/0003-async-pool-janitor-and-awol.md).

## [1.9.5] - 2026-10-16
- Web UI post-release feedback and mobile layout fixes (#152). StarRocks statements run through Flight SQL Execute (#153).

## [1.9.4] - 2026-10-15
- Image rebuilt on patched base images (#151).

## [1.9.3] - 2026-10-13
- ECS polling continues when DescribeTasks omits a starting, pending or provisioning task (#150).
- Web UI moved from react-ui to `@patterninc/pattern-ui` (#149).

## [1.9.2] - 2026-09-29
- Command and cluster API `context` values for `password`, `private_key` and `token` are redacted (#145). Extra names via `sensitive_context_keys`.

## [1.9.1] - 2026-09-29
- One database pool is reused and queries stay on the session transaction (#146).

## [1.9.0] - 2026-09-08
- sparkeks `script_uri` and `entry_point` on the SQL wrapper (#143).

## [1.8.4] - 2026-08-31
- Job file endpoints (`/api/v1/job/<id>/<file>`) are limited to the user who submitted the job (#140).
- **Migration:** callers that fetched another user's job output now get an error.
