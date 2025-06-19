create table if not exists job_command_tags
(
    system_job_id bigint not null,
    command_tag varchar(128) not null,
    constraint _job_command_tags_pk primary key (system_job_id, command_tag),
    constraint _job_command_tags_command_tag unique (command_tag, system_job_id)
);

alter table job_command_tags alter column command_tag type varchar(128);
