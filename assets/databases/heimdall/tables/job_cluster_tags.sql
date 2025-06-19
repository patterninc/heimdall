create table if not exists job_cluster_tags
(
    system_job_id bigint not null,
    cluster_tag varchar(128) not null,
    constraint _job_cluster_tags_pk primary key (system_job_id, cluster_tag),
    constraint _job_cluster_tags_cluster_tag unique (cluster_tag, system_job_id)
);

alter table job_cluster_tags alter column cluster_tag type varchar(128);
