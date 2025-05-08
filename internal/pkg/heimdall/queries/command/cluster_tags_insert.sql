insert into command_cluster_tags
(
    system_command_id,
    command_cluster_tag
)
values
    {{ .Slice }};