# Runbooks

Operational procedures for Heimdall. They are written from the code and config in this repo (2026-11-19). Anything that depends on how Heimdall is
deployed outside this repo (cluster, namespace, dashboards, who to page) is marked `TODO(owner)`; fill it in rather than guessing.

| Runbook | Use when |
| --- | --- |
| [deploy-and-rollback.md](deploy-and-rollback.md) | Shipping a version or rolling one back |
| [stuck-or-awol-jobs.md](stuck-or-awol-jobs.md) | Jobs stay `RUNNING` or fail as AWOL |
| [cluster-health.md](cluster-health.md) | `/api/v1/clusters/health` returns 503 |
| [database-schema-deploy.md](database-schema-deploy.md) | Applying schema changes with `deploydb` |

Incident triage: [../oncall.md](../oncall.md).
