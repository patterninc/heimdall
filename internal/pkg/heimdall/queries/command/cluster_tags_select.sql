select
    command_cluster_tag
from
    command_cluster_tags
where
    system_command_id = $1;
