delete from active_jobs
where
    system_job_id in ( {{ .Slice }} )
;
