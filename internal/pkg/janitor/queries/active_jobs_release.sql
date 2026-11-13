update active_jobs
set
    agent_name = null
where
    system_job_id in ({{ .Slice }})
    and last_heartbeat < (extract(epoch from now())::int - $1)
;
