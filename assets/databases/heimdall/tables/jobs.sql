create table if not exists jobs
(
    system_job_id bigint generated always as identity,
    job_command_id int null,
    job_cluster_id int null,
    job_status_id smallint not null,
    job_id varchar(64) not null,
    job_name varchar(64) not null,
    job_version varchar(32) not null,
    job_description varchar(255) null,
    job_context varchar(65535) null,
    job_error varchar(1024) null,
    username varchar(64) not null,
    is_sync boolean not null,
    store_result_sync boolean not null default false,
    created_at int not null default extract(epoch from now())::int,
    updated_at int not null default extract(epoch from now())::int,
    constraint _jobs_pk primary key (system_job_id),
    constraint _jobs_job_id unique (job_id)
);

alter table jobs add column if not exists store_result_sync boolean not null default false;
alter table jobs add column if not exists canceled_by varchar(64) null;

-- Originally had "cancelled_by" column and "cancelling" status, but we aren't british. Whoops.
do $$ begin
    if exists (select 1 from information_schema.columns where table_name = 'jobs' and column_name = 'cancelled_by') then
        update jobs set canceled_by = cancelled_by where canceled_by is null and cancelled_by is not null;
        alter table jobs drop column cancelled_by;
    end if;
end $$;
update jobs set canceled_by = '' where canceled_by is null;