select
    c.cluster_status_id,
    c.updated_at
from
    clusters c
where
    c.cluster_id = $1;
