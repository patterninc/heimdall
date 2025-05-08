update active_jobs j
set
    last_heartbeat = extract(epoch from now())::int
where
    j.system_job_id = $1 and
    j.agent_name = $2
;
