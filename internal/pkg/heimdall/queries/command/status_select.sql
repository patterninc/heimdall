select
    c.command_status_id,
    c.updated_at
from
    commands c
where
    c.command_id = $1;
