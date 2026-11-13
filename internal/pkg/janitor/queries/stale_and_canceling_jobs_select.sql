with stale_ids as (
  -- stale jobs still in active_jobs table (worker stopped heartbeating); released jobs only once awol.
  -- locking rechecks a row an agent claimed or heartbeated after this query's snapshot
  select aj.system_job_id
  from active_jobs aj
  where aj.last_heartbeat > 0
    and aj.last_heartbeat < (extract(epoch from now())::int - $1)
    and (aj.agent_name is not null or aj.last_heartbeat < (extract(epoch from now())::int - $3))
  for update skip locked
),
candidate_ids as (
  -- canceling jobs that were never claimed by an agent or stopped heartbeating
  select j.system_job_id
  from jobs j
  left join active_jobs aj on aj.system_job_id = j.system_job_id
  where j.job_status_id = 7
    and (aj.system_job_id is null or aj.last_heartbeat is null or aj.agent_name is null)

  union

  select system_job_id from stale_ids
),
picked as (
  select
    j.system_job_id,
    j.job_id,
    j.job_status_id,
    j.job_command_id,
    j.job_cluster_id,
    coalesce(aj.last_heartbeat < (extract(epoch from now())::int - $3), false) as awol
  from jobs j
  join candidate_ids c on c.system_job_id = j.system_job_id
  left join active_jobs aj on aj.system_job_id = j.system_job_id
  order by j.system_job_id
  for update of j skip locked
  limit $2
)
select
  p.system_job_id,
  p.job_id,
  p.job_status_id,
  cm.command_id,
  cl.cluster_id,
  p.awol
from picked p
join commands cm on cm.system_command_id = p.job_command_id
join clusters cl on cl.system_cluster_id = p.job_cluster_id
order by p.system_job_id;