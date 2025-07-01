update clusters
set
    cluster_status_id = $2,
    username = $3,
    updated_at = extract(epoch from now())::int
where
    cluster_id = $1;
