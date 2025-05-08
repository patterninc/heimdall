delete from cluster_tags
where
    system_cluster_id = $1;
