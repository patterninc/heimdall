
DELETE FROM job_command_tags
WHERE system_job_id IN (
  SELECT system_job_id FROM jobs WHERE updated_at < extract(epoch FROM now() - ($1 || ' days')::interval)::int
);
