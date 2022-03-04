package test

import (
	"todo/src/entity"

	"github.com/google/uuid"
)

func GetNoteWithMaxVersion() *entity.Note {
	note := &entity.Note{
		NoteGuid:          new(string),
		Version:           new(int8),
		Text:              new(string),
		PrevNoteVersionId: new(int64),
		UserId:            new(int64),
		Deleted:           new(bool),
		Archive:           new(bool),
	}
	*note.NoteGuid = "not guid1"
	*note.Version = 1
	*note.Text = "note text1_2"
	*note.PrevNoteVersionId = 1
	*note.UserId = 1
	*note.Deleted = false
	*note.Archive = false
	return note
}

func CreateNewRandomNote() *entity.Note {
	noteFiles := []entity.NoteFile{
		*createNewRandomNoteFile(),
		*createNewRandomNoteFile(),
		*createNewRandomNoteFile(),
	}

	note := &entity.Note{
		NoteGuid:  new(string),
		Version:   new(int8),
		Text:      new(string),
		UserId:    new(int64),
		Deleted:   new(bool),
		Archive:   new(bool),
		NoteFiles: noteFiles,
	}
	randomUuid := uuid.New().String()
	*note.NoteGuid = randomUuid
	*note.Text = randomUuid
	*note.UserId = 1
	return note
}

func createNewRandomNoteFile() *entity.NoteFile {
	randomUuid := uuid.New().String()
	file := &entity.NoteFile{
		Filename: &randomUuid,
		Data:     []byte{},
	}
	return file
}
