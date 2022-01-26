create table note (
    id          int not null auto_increment primary key,
    note_guid   varchar(16),
    version     int not null default 0,
    text        text,
    user_id     int not null,
    create_date datetime DEFAULT now(),
    deleted TINYINT(1) default 0,
    archive TINYINT(1) default 0
);

