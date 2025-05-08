insert into command_tags
(
    system_command_id,
    command_tag
)
values
    {{ .Slice }};