# 0003. Async jobs run in a worker pool; a janitor releases stale jobs and fails them as AWOL

- Status: Accepted (retroactive, 2026-11-19; "AWOL plan" phases 1-4, PRs #161, #162, #164)
- Evidence: `cmd/heimdall/heimdall.go` (defaults), `internal/pkg/pool`, `internal/pkg/janitor`, `internal/pkg/heimdall/jobs_async.go`, `internal/pkg/heimdall/keepalive.go`

## Decision
- An agent claims async jobs into a worker `pool` (`pool.size`) and keeps a heartbeat per claimed job (`janitor.keepalive`, default 5 s).
- On SIGTERM the agent drains in-flight jobs for up to `HEIMDALL_SHUTDOWN_TIMEOUT` seconds (default 30) before exiting.
- The janitor releases a resumable job whose heartbeat is older than `stale_job` (default 45 s) so another agent can resume it (sparkeks and ecs). A released job nobody picks up is failed as AWOL after `awol_timeout` (default 3600 s). Non-resumable jobs are failed instead of released.

## Consequences
- A pod restart or deploy does not fail resumable jobs; it hands them to another agent.
- `stale_job` must stay above the keepalive interval with margin, or healthy jobs are released.
- Behaviour is covered by `internal/pkg/janitor/stale_jobs_test.go`, which needs `HEIMDALL_TEST_DATABASE`.
