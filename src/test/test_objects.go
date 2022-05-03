package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"todo/src/entity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetNoteWithMaxVersion() *entity.Note {
	note := &entity.Note{
		NoteGuid:          new(string),
		Version:           new(int8),
		Title:             new(string),
		Text:              new(string),
		PrevNoteVersionId: new(int64),
		UserId:            new(int64),
		Deleted:           new(bool),
		Archive:           new(bool),
	}
	*note.NoteGuid = "not guid1"
	*note.Version = 1
	*note.Title = "note title1_2"
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
		Title:      new(string),
		Text:      new(string),
		UserId:    new(int64),
		Deleted:   new(bool),
		Archive:   new(bool),
		NoteFiles: noteFiles,
	}
	randomUuid := uuid.New().String()
	*note.NoteGuid = randomUuid
	*note.Title = randomUuid
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

func ParseResponseBody(t *testing.T, res *http.Response) gin.H {
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("ошибка чтения response body: %s", err)
	}
	var got gin.H
	err = json.Unmarshal(bodyBytes, &got)
	if err != nil {
		t.Fatalf("ошибка формирования json: %s", err)
	}
	return got
}

func CreateAndGetNewNoteWithNVersion(t *testing.T, noteService entity.NoteService, versionNumber int) *entity.Note {
	if (versionNumber < 1) {
		t.Fatalf("нельзя создать note с 0 версий")
	}

	/* Создаём новый note */
	note := CreateNewRandomNote()
	note.Title = new(string)
	*note.Title = strconv.Itoa(1)
	note.Text = new(string)
	*note.Text = strconv.Itoa(1)
	noteId, err := noteService.SaveNote(note)
	if err != nil {
		t.Fatalf("error was not expected while test method: %s", err)
	}
	note, err = noteService.GetNote(*noteId)
	if err != nil {
		t.Fatalf("error was not expected while test method: %s", err)
	}

	for i:=1; i < int(versionNumber); i++ {
		*note.Text = strconv.Itoa(i + 1)
		/* Создаём новую версию note */
		noteId, err = noteService.SaveNote(note)
		if err != nil {
			t.Fatalf("error was not expected while test method: %s", err)
		}
		note, err = noteService.GetNote(*noteId)
		if err != nil {
			t.Fatalf("error was not expected while test method: %s", err)
		}
	}
	
	return note
}