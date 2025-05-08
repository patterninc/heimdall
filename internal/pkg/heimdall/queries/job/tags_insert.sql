insert into job_tags
(
    system_job_id,
    job_tag
)
values
    {{ .Slice }};