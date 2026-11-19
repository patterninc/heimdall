# Stuck or AWOL jobs

Background: [ADR 0003](../adr/0003-async-pool-janitor-and-awol.md). Defaults: `stale_job` 45 s, `awol_timeout` 3600 s, janitor `clean_interval` 60 s, keepalive 5 s.

1. Find the job: `GET /api/v1/job/<id>` and `GET /api/v1/job/<id>/status`. Fetch `GET /api/v1/job/<id>/stderr` and `/stdout` for plugin output.
2. **Running, no progress.** If the agent that owns it is gone, the janitor releases a resumable job after `stale_job`; another agent resumes it. If it is not resumed, check that at least one agent is running and its pool is not full (`pool.size`).
3. **Failed as AWOL.** The job was released and nobody resumed it within `awol_timeout`. Check agent logs around the release time for crashes or a drain that exceeded `HEIMDALL_SHUTDOWN_TIMEOUT`. Resubmit with `POST /api/v1/job`.
4. **Cancel.** `POST /api/v1/job/<id>/cancel` (async jobs only).
5. **Performance.** Set `ENABLE_DEBUG=true` to expose `/debug/pprof/`; unset it afterwards.
6. Metrics: the server proxies `/metrics` to the Prometheus listener at `PROMETHEUS_ADDRESS` when `METRICS_EMITTER=prometheus`.
