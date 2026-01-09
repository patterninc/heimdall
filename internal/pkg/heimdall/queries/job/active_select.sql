with a as
(
    select
        jj.system_job_id,
        cc.command_id,
        cl.cluster_id,
        jj.job_status_id,
        jj.job_id,
        jj.job_name,
        jj.job_version,
        jj.job_description,
        jj.job_context,
        jj.username,
        jj.is_sync,
        jj.created_at,
        jj.updated_at,
        jj.store_result_sync
    from
        active_jobs aj
        join jobs jj on jj.system_job_id = aj.system_job_id
        join commands cc on cc.system_command_id = jj.job_command_id
        join clusters cl on cl.system_cluster_id = jj.job_cluster_id
    where
        aj.agent_name is null
        and jj.job_status_id != 7  -- Not canceling. Jobs can be canceled before being assigned to an agent.
    order by
        aj.system_job_id
    for update
        skip locked
    limit
        $1
)
update active_jobs j
set
    agent_name = $2,
    last_heartbeat = extract(epoch from now())::int
from
    a
where
    j.system_job_id = a.system_job_id
returning
    a.system_job_id,
    a.command_id,
    a.cluster_id,
    a.job_status_id,
    a.job_id,
    a.job_name,
    a.job_version,
    a.job_description,
    a.job_context,
    a.username,
    a.is_sync,
    a.created_at,
    a.updated_at,
    a.store_result_sync;
