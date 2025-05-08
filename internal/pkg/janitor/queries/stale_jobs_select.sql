select
    system_job_id
from
    active_jobs
where
    last_heartbeat > 0 and
    extract(epoch from now())::int - $1 > last_heartbeat
limit
    25
;
