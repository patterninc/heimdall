update commands
set
    command_status_id = $2,
    username = $3,
    updated_at = extract(epoch from now())::int
where
    command_id = $1;
