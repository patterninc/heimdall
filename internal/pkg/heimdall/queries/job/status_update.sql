update jobs
set
    job_status_id = $1,
    job_error = left($2::text, 1024),
    extra_job_attributes = coalesce(nullif($3::text, '')::jsonb, extra_job_attributes),
    updated_at = extract(epoch from now())::int
where
    system_job_id = $4;
