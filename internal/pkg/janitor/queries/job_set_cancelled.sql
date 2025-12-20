update jobs
set
    job_status_id = 8,  -- CANCELLED
    updated_at = extract(epoch from now())::int
where
    system_job_id in ( {{ .Slice }} )
;

