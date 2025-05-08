create table if not exists job_statuses
(
    job_status_id smallint not null,
    job_status_name varchar(32) not null,
    constraint _job_statuses_pk primary key (job_status_id),
    constraint _job_statuses_job_status_name unique (job_status_name)
);
