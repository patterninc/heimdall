update jobs
set
    job_status_id = (select job_status_id from job_statuses where job_status_name = 'CANCELLED'),
    updated_at = extract(epoch from now())::int
where
    job_id = $1;
