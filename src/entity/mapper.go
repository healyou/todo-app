package entity

import (
	"database/sql"

	"github.com/pkg/errors"
)

func MapNote(noteRow *sql.Row) (*Note, error) {
	return mapOneNote(noteRow, nil)
}

func mapOneNote(noteRow *sql.Row, noteRows *sql.Rows) (*Note, error) {
	var note Note

	var deleted int8
	var archive int8
	var actual int8
	if (noteRow != nil) {
		err := noteRow.Scan(
			&note.Id, &note.NoteGuid, &note.Version, &note.PrevNoteVersionId, &note.Text, &note.UserId, &note.CreateDate, 
			&deleted, &archive, &actual)
		if err != nil {
			return nil, errors.Wrap(err, "ошибка маппинга note")
		}
	} else {
		err := noteRows.Scan(
			&note.Id, &note.NoteGuid, &note.Version, &note.PrevNoteVersionId, &note.Text, &note.UserId, &note.CreateDate, 
			&deleted, &archive, &actual)
		if err != nil {
			return nil, errors.Wrap(err, "ошибка маппинга note")
		}
	}
	note.Deleted = new(bool)
	note.Archive = new(bool)
	note.Actual = new(bool)
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
	if actual == 1 {
		*note.Actual = true
	} else {
		*note.Actual = false
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
			return nil, errors.Wrap(err, "ошибка маппинга notefile")
		}
		noteFiles = append(noteFiles, file)
	}

	return noteFiles, nil
}

func MapNotes(noteRows *sql.Rows) ([]Note, error) {
	var notes []Note

	for noteRows.Next() {
		note, err := mapOneNote(nil, noteRows)
		if err != nil {
			return nil, errors.Wrap(err, "ошибка маппинга notes")
		}
		notes = append(notes, *note)
	}

	return notes, nil
}

func MapNoteVersionHistoryItems(noteVersionInfoRows *sql.Rows) ([]NoteVersionInfo, error) {
	var noteVersionHistoryItems []NoteVersionInfo

	for noteVersionInfoRows.Next() {
		item, err := MapOneNoteVersionHistory(noteVersionInfoRows)
		if err != nil {
			return nil, errors.Wrap(err, "ошибка маппинга NoteVersionInfo")
		}
		noteVersionHistoryItems = append(noteVersionHistoryItems, *item)
	}

	return noteVersionHistoryItems, nil
}

func MapOneNoteVersionHistory(noteVersionHistoryRows *sql.Rows) (*NoteVersionInfo, error) {
	var versionInfo NoteVersionInfo

	var actual int8
	err := noteVersionHistoryRows.Scan(
		&versionInfo.NoteId, &versionInfo.PrevNoteVersionId, &versionInfo.Version, &versionInfo.CreateDate, &actual)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка маппинга NoteVersionInfo")
	}
	versionInfo.Actual = new(bool)
	if actual == 1 {
		*versionInfo.Actual = true
	} else {
		*versionInfo.Actual = false
	}

	return &versionInfo, nil
}