select
    c.system_command_id,
    c.command_status_id,
    c.command_id,
    c.command_name,
    c.command_version,
    c.command_plugin,
    c.command_description,
    c.command_context,
    c.username,
    c.is_sync,
    c.created_at,
    c.updated_at
from
    commands c
    join command_statuses cs on cs.command_status_id = c.command_status_id{{ if .Clause }}
where
    {{ .Clause }}{{end}}
order by
    c.system_command_id desc
limit
    101
;
