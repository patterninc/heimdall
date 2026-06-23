update jobs
set
    job_status_id = $1,
    job_error = left($2::text, 1024),
    spark_application_id = coalesce(nullif($3::text, ''), spark_application_id),
    updated_at = extract(epoch from now())::int
where
    system_job_id = $4;
