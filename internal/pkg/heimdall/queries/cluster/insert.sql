insert into clusters
(
    cluster_status_id,
    cluster_id,
    cluster_name,
    cluster_version,
    cluster_description,
    cluster_context,
    username
)
values
($1, $2, $3, $4, $5, $6, $7)
on conflict (cluster_id) do update
set
    cluster_status_id = excluded.cluster_status_id,
    cluster_id = excluded.cluster_id,
    cluster_name = excluded.cluster_name,
    cluster_description = excluded.cluster_description,
    cluster_context = excluded.cluster_context,
    username = excluded.username,
    updated_at = extract(epoch from now())::int
returning system_cluster_id;
