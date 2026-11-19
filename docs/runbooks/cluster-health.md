# Cluster health check failing

`GET /api/v1/clusters/health` returns 200 when every opted-in cluster (`health_check: true`) passes and 503 when any fails; each check reports `ok`, `error` or `unchecked`.
`GET /api/v1/cluster/<id>/health` probes one cluster. The probe timeout is `health_check.timeout_seconds` (default 30 s in README; `configs/local.yaml` sets 300).

1. Identify the failing cluster and its plugin from the response.
2. Reproduce the plugin's lightweight probe from the agent host (for example Glue checks the configured `catalog_id`; see PR #142).
3. Check credentials and network reachability for that backend.
4. A failing or misconfigured health check does not affect job execution or startup (README, "Cluster Health Checks"). To take a cluster out of rotation while it is fixed, set its status with `PUT /api/v1/cluster/<id>/status`.
