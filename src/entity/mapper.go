package entity

import (
	"database/sql"
)

func MapNote(noteRow *sql.Row) (*Note, error) {
	var note Note

	var deleted int8
	var archive int8
	err := noteRow.Scan(
		&note.Id, &note.NoteGuid, &note.Version, &note.Text, &note.UserId, &note.CreateDate, &deleted,
		&archive)
	if err != nil {
		return nil, err
	}
	note.Deleted = new(bool)
	note.Archive = new(bool)
	if deleted == 1 {
		*note.Deleted = true
	} else {
		*note.Deleted = false
	}
	if archive == 1 {
		*note.Archive = true
	} else {
		*note.Archive = false
	}

	return &note, nil
}

func MapNoteFiles(noteFileRows *sql.Rows) ([]NoteFile, error) {
	var noteFiles []NoteFile

	for noteFileRows.Next() {
		var file NoteFile
		err := noteFileRows.Scan(
			&file.Id, &file.NoteId, &file.Guid, &file.Filename)
		if err != nil {
			return nil, err
		}
		noteFiles = append(noteFiles, file)
	}

	return noteFiles, nil
}
