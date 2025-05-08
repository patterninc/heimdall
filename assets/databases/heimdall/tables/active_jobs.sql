create table if not exists active_jobs
(
    system_job_id bigint not null,
    agent_name varchar(64) null,
    last_heartbeat int null,
    constraint _active_jobs_pk primary key (system_job_id)
);
