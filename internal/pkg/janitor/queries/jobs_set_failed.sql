update jobs
set
    job_status_id = 4, -- failed
    job_error = 'awol: job is stale, marking failed',
    updated_at = extract(epoch from now())::int
where
    system_job_id in ( {{ .Slice }} ) and
    job_status_id = 3 -- running
;
