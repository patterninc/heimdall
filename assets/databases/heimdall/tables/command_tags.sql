create table if not exists command_tags
(
    system_command_id int not null,
    command_tag varchar(128) not null,
    constraint _command_tags_pk primary key (system_command_id, command_tag),
    constraint _commands_tags_command_tag unique (command_tag, system_command_id)
);

alter table command_tags alter column command_tag type varchar(128);
