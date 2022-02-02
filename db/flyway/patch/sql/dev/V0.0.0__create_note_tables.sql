create table note (
    id          int not null auto_increment primary key,
    note_guid   varchar(36),
    version     int not null default 0,
    text        text,
    user_id     int not null,
    create_date datetime DEFAULT now(),
    deleted TINYINT(1) default 0,
    archive TINYINT(1) default 0,
    UNIQUE KEY unique_user_version_note_index (note_guid,version,user_id)
);
