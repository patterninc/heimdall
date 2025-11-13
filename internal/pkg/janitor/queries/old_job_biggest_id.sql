SELECT system_job_id 
FROM jobs 
WHERE updated_at < $1
ORDER BY updated_at desc 
LIMIT 1