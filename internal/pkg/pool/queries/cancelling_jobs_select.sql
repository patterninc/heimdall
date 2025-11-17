select
    j.job_id
from
    jobs j
    join job_statuses js on j.job_status_id = js.job_status_id
where
    js.job_status_name = 'CANCELLING'
limit
    25;
