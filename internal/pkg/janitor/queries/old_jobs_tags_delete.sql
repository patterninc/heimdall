DELETE FROM job_tags
WHERE system_job_id IN (
  SELECT system_job_id
  FROM job_tags
  WHERE system_job_id <= $1
  LIMIT 1000
);