select
    command_tag
from
    job_command_tags
where
    system_job_id = $1;
