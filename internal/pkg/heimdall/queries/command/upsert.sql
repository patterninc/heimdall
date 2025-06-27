insert into commands
(
    command_status_id,
    command_id,
    command_name,
    command_version,
    command_plugin,
    command_description,
    command_context,
    username,
    is_sync
)
values
($1, $2, $3, $4, $5, $6, $7, $8, $9)
on conflict (command_id) do update
set
    command_status_id = excluded.command_status_id,
    command_id = excluded.command_id,
    command_name = excluded.command_name,
    command_version = excluded.command_version,
    command_plugin = excluded.command_plugin,
    command_description = excluded.command_description,
    command_context = excluded.command_context,
    username = excluded.username,
    is_sync = excluded.is_sync,
    updated_at = extract(epoch from now())::int
returning system_command_id;