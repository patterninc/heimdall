# Apply schema changes (`deploydb`)

Background: [ADR 0002](../adr/0002-postgres-state-and-deploydb.md).

1. Add or edit the SQL file under `assets/databases/heimdall/tables/` (or `data/`) and list it in `assets/databases/heimdall/build/heimdall.lst`. SQL must be re-runnable.
2. Build: `./build.sh --go` (produces `dist/deploydb`).
3. Apply: `./dist/deploydb -conf <heimdall.yaml> assets/databases/heimdall/build/heimdall.lst`. The conf needs `database.connection_string`.
4. Verify with the integration test: `HEIMDALL_TEST_DATABASE=<connection string> go test ./internal/pkg/janitor/...`.
5. Deploy the server version that uses the change.

Local: `docker compose up postgres` exposes Postgres on `localhost:5433` (user/db `heimdall`).
