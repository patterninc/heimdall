create table if not exists command_statuses
(
    command_status_id smallint not null,
    command_status_name varchar(32) not null,
    constraint _command_statuses_pk primary key (command_status_id),
    constraint _command_statuses_command_status_name unique (command_status_name)
);
