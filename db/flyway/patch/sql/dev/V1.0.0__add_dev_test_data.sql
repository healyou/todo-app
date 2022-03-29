INSERT INTO note (note_guid, version, text, actual, user_id)
VALUES ('not guid1', 0, 'note text1_1', 0, 1);
INSERT INTO note (note_guid, version, prev_note_version_id, text, actual, user_id)
VALUES ('not guid1', 1, 1, 'note text1_2', 1, 1);
INSERT INTO note (note_guid, text, user_id)
VALUES ('not guid2', 'note text2', 2);
INSERT INTO note (note_guid, text, user_id)
VALUES ('not guid3', 'note text3', 3);