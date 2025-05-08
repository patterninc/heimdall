create table if not exists commands
(
    system_command_id int generated always as identity,
    command_status_id smallint not null,
    command_id varchar(32) not null,
    command_name varchar(32) not null,
    command_version varchar(32) not null,
    command_plugin varchar(64) not null,
    command_description varchar(255) null,
    command_context varchar(65535) null,
    username varchar(64) not null,
    is_sync boolean not null,
    created_at int not null default extract(epoch from now())::int,
    updated_at int not null default extract(epoch from now())::int,
    constraint _commands_pk primary key (system_command_id),
    constraint _commands_command_id unique (command_id)
);
