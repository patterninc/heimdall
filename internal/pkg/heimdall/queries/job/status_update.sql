update jobs
set
    job_status_id = $1,
    job_error = case when $2 is null then null else left($2, 1024) end,
    updated_at = extract(epoch from now())::int
where
    system_job_id = $3;
