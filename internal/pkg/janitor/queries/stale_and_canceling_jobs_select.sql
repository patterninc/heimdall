select
    j.system_job_id,
    j.job_id,
    j.job_status_id,
    cm.command_id,
    cl.cluster_id
from
    jobs j
    left join active_jobs aj on j.system_job_id = aj.system_job_id
    join commands cm on cm.system_command_id = j.job_command_id
    join clusters cl on cl.system_cluster_id = j.job_cluster_id
where
    (
        -- Stale jobs: must be in active_jobs and have heartbeat timeout
        (aj.system_job_id is not null and aj.last_heartbeat > 0 and extract(epoch from now())::int - $1 > aj.last_heartbeat)
        or
        -- Canceling jobs: status is CANCELING
        j.job_status_id = 7
    )
order by
    j.system_job_id
for update of j
    skip locked
limit
    $2
;

