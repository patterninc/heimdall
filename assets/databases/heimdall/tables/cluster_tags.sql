create table if not exists cluster_tags
(
    system_cluster_id int not null,
    cluster_tag varchar(128) not null,
    constraint _cluster_tags_pk primary key (system_cluster_id, cluster_tag),
    constraint _cluster_tags_cluster_tag unique (cluster_tag, system_cluster_id)
);

alter table cluster_tags alter column cluster_tag type varchar(128);
