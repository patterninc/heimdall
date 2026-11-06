update active_jobs
set
    agent_name = null,
    last_heartbeat = null
where
    system_job_id in ({{ .Slice }})
;
