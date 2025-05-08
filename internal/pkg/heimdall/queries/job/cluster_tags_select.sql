select
    cluster_tag
from
    job_cluster_tags
where
    system_job_id = $1;
