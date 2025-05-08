insert into cluster_statuses
(
    cluster_status_id,
    cluster_status_name
)
values
    (1, 'ACTIVE'),
    (2, 'INACTIVE'),
    (3, 'DELETED')
on conflict (cluster_status_id) do update
set
    cluster_status_name = excluded.cluster_status_name;
