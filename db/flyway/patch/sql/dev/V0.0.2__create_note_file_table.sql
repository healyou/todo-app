create table note_file (
    id int not null auto_increment primary key,
    note_id int not null,
    file_guid varchar(36) not null,
    filename varchar(256) not null,
    FOREIGN KEY(note_id) REFERENCES note(id)
);