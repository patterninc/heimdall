SELECT system_job_id 
FROM jobs 
WHERE updated_at < $1
ORDER BY system_job_id desc 
LIMIT 1