select
    j.system_job_id,
    j.job_status_id,
    j.job_name,
    j.job_version,
    j.job_description,
    j.job_context,
    j.job_error,
    j.username,
    j.is_sync,
    j.created_at,
    j.updated_at,
    cm.command_id,
    cm.command_name,
    cl.cluster_id,
    cl.cluster_name
from
    jobs j
    left join commands cm on cm.system_command_id = j.job_command_id
    left join clusters cl on cl.system_cluster_id = j.job_cluster_id
where
    j.job_id = $1;
