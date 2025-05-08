insert into job_command_tags
(
    system_job_id,
    command_tag
)
values
    {{ .Slice }};