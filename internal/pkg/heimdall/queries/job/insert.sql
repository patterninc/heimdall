insert into jobs
(
    job_command_id,
    job_cluster_id,
    job_status_id,
    job_id,
    job_name,
    job_version,
    job_description,
    job_context,
    job_error,
    username,
    is_sync,
    store_result_sync,
    canceled_by
)
select
    cm.system_command_id,
    cl.system_cluster_id,
    $3, -- job_status_id
    $4, -- job_id
    $5, -- job_name
    $6, -- job_version
    $7, -- job_description
    $8, -- job_context
    case when $9 is null then null else left($9, 1024) end, -- job_error
    $10, -- username
    $11, -- is_sync
    $12, -- store_result_sync
    $13 -- canceled_by
from
    clusters cl,
    commands cm
where
    cl.cluster_id = $1 and
    cm.command_id = $2
returning system_job_id;
