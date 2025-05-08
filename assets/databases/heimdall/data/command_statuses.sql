insert into command_statuses
(
    command_status_id,
    command_status_name
)
values
    (1, 'ACTIVE'),
    (2, 'INACTIVE'),
    (3, 'DELETED')
on conflict (command_status_id) do update
set
    command_status_name = excluded.command_status_name;
