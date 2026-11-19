# On-call playbook

Heimdall is the job-execution service the data platform submits work through. This playbook covers what can be checked from this repo. Paging, dashboards and
deployment details are not recorded here: `TODO(owner)`.

| Symptom | First checks | Runbook |
| --- | --- | --- |
| API returns 5xx or is unreachable | Agent process/pod up? Database reachable (`database.connection_string`)? | [deploy-and-rollback](runbooks/deploy-and-rollback.md) |
| Jobs stay `RUNNING` or become AWOL | Agents running, pool not full, heartbeat | [stuck-or-awol-jobs](runbooks/stuck-or-awol-jobs.md) |
| One backend fails, others fine | `GET /api/v1/clusters/health`, plugin stderr | [cluster-health](runbooks/cluster-health.md) |
| Schema error after a deploy | Was `deploydb` run for this release? | [database-schema-deploy](runbooks/database-schema-deploy.md) |

Triage order: (1) is the server answering `GET /api/v1/clusters/health`, (2) is Postgres up, (3) are agents claiming jobs (`GET /api/v1/jobs`), (4) is one plugin or cluster the common factor.
Escalation contacts: `TODO(owner)` (no CODEOWNERS or Slack channel is recorded in this repo yet).
