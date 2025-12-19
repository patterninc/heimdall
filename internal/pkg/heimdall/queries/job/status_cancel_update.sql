update jobs
set
    job_status_id = 7,  -- CANCELLING
    cancelled_by = $2,
    updated_at = extract(epoch from now())::int
where
    job_id = $1 
    and job_status_id not in (4, 5, 6, 8);  -- Not in FAILED, KILLED, SUCCEEDED, CANCELLED
