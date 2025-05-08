insert into job_cluster_tags
(
    system_job_id,
    cluster_tag
)
values
    {{ .Slice }};