INSERT INTO note (note_guid, text, user_id)
VALUES ('not guid1', 'note text1', 1);
INSERT INTO note (note_guid, text, user_id)
VALUES ('not guid2', 'note text2', 2);
INSERT INTO note (note_guid, text, user_id)
VALUES ('not guid3', 'note text3', 3);

INSERT INTO note_file (id, note_id, file_guid, filename)
VALUES (1, 1, 'file_guid', 'filename');
INSERT INTO note_file (id, note_id, file_guid, filename)
VALUES (2, 2, 'file_guid', 'filename');
INSERT INTO note_file (id, note_id, file_guid, filename)
VALUES (3, 3, 'file_guid', 'filename');