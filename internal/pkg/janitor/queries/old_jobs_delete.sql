
DELETE FROM jobs
WHERE system_job_id IN (
  SELECT system_job_id
  FROM jobs
  WHERE system_job_id <= $1
  LIMIT 1000
);