create table if not exists clusters
(
    system_cluster_id int generated always as identity,
    cluster_status_id smallint not null,
    cluster_id varchar(32) not null,
    cluster_name varchar(32) not null,
    cluster_version varchar(32) not null,
    cluster_description varchar(255) null,
    cluster_context varchar(65535) null,
    username varchar(64) not null,
    created_at int not null default extract(epoch from now())::int,
    updated_at int not null default extract(epoch from now())::int,
    constraint _clusters_pk primary key (system_cluster_id),
    constraint _clusters_cluster_id unique (cluster_id)
);
