create table if not exists command_cluster_tags
(
    system_command_id int not null,
    command_cluster_tag varchar(128) not null,
    constraint _command_cluster_tags_pk primary key (system_command_id, command_cluster_tag)
);

alter table command_cluster_tags alter column command_cluster_tag type varchar(128);
