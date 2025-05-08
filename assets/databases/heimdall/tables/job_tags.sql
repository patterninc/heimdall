create table if not exists job_tags
(
    system_job_id bigint not null,
    job_tag varchar(64) not null,
    constraint _job_tags_pk primary key (system_job_id, job_tag),
    constraint _job_tags_job_tag unique (job_tag, system_job_id)
);
