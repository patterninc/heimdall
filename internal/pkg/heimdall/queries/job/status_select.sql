select
    j.job_status_id,
    j.job_error,
    j.updated_at
from
    jobs j
where
    j.job_id = $1;
