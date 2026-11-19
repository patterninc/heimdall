# 0002. Postgres is the only state store; schema applied by `deploydb`

- Status: Accepted (retroactive, 2026-11-19)
- Evidence: `assets/databases/heimdall/build/heimdall.lst`, `cmd/deploydb/deploydb.go`, `internal/pkg/database/deploy.go`, `docker-compose.yaml`, `configs/local.yaml` (`database.connection_string`)

## Decision
Jobs, commands, clusters, tags and status lookups live in PostgreSQL. The schema is a set of idempotent SQL files listed, in order, in
`assets/databases/heimdall/build/heimdall.lst` (tables first, then static status data). `cmd/deploydb` reads that list and applies each file.
There is no versioned migration tool; the list file is the only ordering mechanism.

## Consequences
- A schema change is a new or edited SQL file plus an entry in `heimdall.lst`, and `deploydb` must be run before the new server version starts.
- SQL files must be safe to re-run. (inferred from how `deploy.go` applies every file on every run)
- Database tests that need the schema read `HEIMDALL_TEST_DATABASE` and skip when it is unset (`internal/pkg/janitor/stale_jobs_test.go`).
