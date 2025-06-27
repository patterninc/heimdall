select
    c.system_cluster_id,
    c.cluster_status_id,
    c.cluster_name,
    c.cluster_version,
    c.cluster_description,
    c.cluster_context,
    c.username,
    c.created_at,
    c.updated_at
from
    clusters c
where
    c.cluster_id = $1;
