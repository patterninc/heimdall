select
    job_tag
from
    job_tags
where
    system_job_id = $1;
