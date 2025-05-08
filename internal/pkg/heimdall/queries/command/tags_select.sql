select
    command_tag
from
    command_tags
where
    system_command_id = $1;
