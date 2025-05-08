create table if not exists cluster_statuses
(
    cluster_status_id smallint not null,
    cluster_status_name varchar(32) not null,
    constraint _cluster_statuses_pk primary key (cluster_status_id),
    constraint _cluster_statuses_cluster_status_name unique (cluster_status_name)
);
