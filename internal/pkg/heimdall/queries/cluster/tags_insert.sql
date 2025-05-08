insert into cluster_tags
(
    system_cluster_id,
    cluster_tag
)
values
    {{ .Slice }};