select
    j.system_job_id,
    j.job_id,
    j.cancellation_ctx,
    cm.command_id,
    cl.cluster_id
from
    jobs j
    join commands cm on cm.system_command_id = j.job_command_id
    join clusters cl on cl.system_cluster_id = j.job_cluster_id
where
    j.job_status_id = 7  -- CANCELLING
order by
    j.system_job_id
for update
    skip locked
limit $1
;

