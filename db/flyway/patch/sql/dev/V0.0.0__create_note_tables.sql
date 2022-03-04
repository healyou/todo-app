create table note (
    id          int not null auto_increment primary key,
    note_guid   varchar(36) not null,
    version     int not null default 0,
    -- Идентификатор записи, с которой была создана текущая запись
    prev_note_version_id int,
    text        text,
    user_id     int not null,
    create_date datetime not null DEFAULT now(),
    deleted TINYINT(1) not null default 0,
    archive TINYINT(1) not null default 0,
    actual TINYINT(1) not null default 1,
    FOREIGN KEY(prev_note_version_id) REFERENCES note(id),
    UNIQUE KEY unique_user_version_note_index (note_guid,version)
);