with candidate_ids as (
  -- canceling jobs
  select j.system_job_id
  from jobs j
  where j.job_status_id = 7

  union

  -- stale active jobs (drive from active_jobs)
  select aj.system_job_id
  from active_jobs aj
  where aj.last_heartbeat > 0
    and aj.last_heartbeat < (extract(epoch from now())::int - $1)
),
picked as (
  select
    j.system_job_id,
    j.job_id,
    j.job_status_id,
    j.job_command_id,
    j.job_cluster_id
  from jobs j
  join candidate_ids c on c.system_job_id = j.system_job_id
  order by j.system_job_id
  for update of j skip locked
  limit $2
)
select
  p.system_job_id,
  p.job_id,
  p.job_status_id,
  cm.command_id,
  cl.cluster_id
from picked p
join commands cm on cm.system_command_id = p.job_command_id
join clusters cl on cl.system_cluster_id = p.job_cluster_id
order by p.system_job_id;

