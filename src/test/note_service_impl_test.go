package test

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
	"todo/src/entity"
)

/* Тест */
func TestGgWp(t *testing.T) {
	txFunc := func(testJdbcTemplate JdbcTemplateImplTest) {
		var noteService = entity.NoteServiceImpl{JdbcTemplate: &testJdbcTemplate}
		err := noteService.Test()
		if err != nil {
			t.Errorf("error was not expected while test method: %s", err)
		}
	}
	ExecuteTestRollbackTransaction(t, txFunc)
}

func TestSaveNewNoteWithData(t *testing.T) {
	txFunc := func(testJdbcTemplate JdbcTemplateImplTest) {
		savedNote := CreateNewRandomNote()

		var noteService = entity.NoteServiceImpl{JdbcTemplate: &testJdbcTemplate}
		result, err := noteService.SaveNote(savedNote)
		if err != nil {
			t.Fatalf("error was not expected while test method: %s", err)
		}

		assert.NotNil(t, result)

		createdNote, err := noteService.GetNote(*result)
		if err != nil {
			t.Fatalf("error was not expected while test method: %s", err)
		}
		assert.Equal(t, *createdNote.Id, *result)
		assert.Equal(t, *createdNote.NoteGuid, *savedNote.NoteGuid)
		assert.NotNil(t, createdNote.Version)
		assert.Equal(t, *createdNote.Text, *savedNote.Text)
		assert.Equal(t, *createdNote.UserId, *savedNote.UserId)
		assert.NotNil(t, createdNote.CreateDate)
		assert.Equal(t, *createdNote.Deleted, false)
		assert.Equal(t, *createdNote.Archive, false)
		assert.NotNil(t, createdNote.NoteFiles)

		assert.NotNil(t, savedNote.NoteFiles)
		assert.Equal(t, len(savedNote.NoteFiles), len(createdNote.NoteFiles))
		for i := 0; i < len(savedNote.NoteFiles); i++ {
			createdFile := createdNote.NoteFiles[i]
			expectedFile := savedNote.NoteFiles[i]
			assert.NotNil(t, createdFile.Id)
			assert.NotNil(t, createdFile.Guid)
			assert.Equal(t, *createdFile.NoteId, *createdNote.Id)
			assert.Equal(t, *createdFile.Filename, *expectedFile.Filename)
		}
	}
	ExecuteTestRollbackTransaction(t, txFunc)
}

func TestGetNote(t *testing.T) {
	txFunc := func(testJdbcTemplate JdbcTemplateImplTest) {
		var noteService = entity.NoteServiceImpl{JdbcTemplate: &testJdbcTemplate}
		result, err := noteService.GetNote(1)
		if err != nil {
			t.Fatalf("error was not expected while test method: %s", err)
		}

		assert.Equal(t, *result.Id, int64(1))
		assert.NotNil(t, *result.NoteGuid)
		assert.NotNil(t, *result.Version)
		assert.NotNil(t, *result.Text)
		assert.NotNil(t, *result.UserId)
		assert.NotNil(t, *result.CreateDate)
		assert.NotNil(t, *result.Deleted)
		assert.NotNil(t, *result.Archive)
		assert.NotNil(t, result.NoteFiles)
	}
	ExecuteTestRollbackTransaction(t, txFunc)
}

func TestGetNoteByNoteGuid(t *testing.T) {
	txFunc := func(testJdbcTemplate JdbcTemplateImplTest) {
		expectedNote := GetNoteWithMaxVersion()

		var noteService = entity.NoteServiceImpl{JdbcTemplate: &testJdbcTemplate}
		result, err := noteService.GetNoteByGuid(*expectedNote.NoteGuid)
		if err != nil {
			t.Fatalf("error was not expected while test method: %s", err)
		}

		assert.NotNil(t, *result.Id)
		assert.Equal(t, *result.NoteGuid, *result.NoteGuid)
		assert.Equal(t, *result.Version, *result.Version)
		assert.Equal(t, *result.Text, *result.Text)
		assert.Equal(t, *result.UserId, *result.UserId)
		assert.NotNil(t, *result.CreateDate)
		assert.Equal(t, *result.Deleted, *result.Deleted)
		assert.Equal(t, *result.Archive, *result.Archive)
		assert.NotNil(t, result.NoteFiles)
	}
	ExecuteTestRollbackTransaction(t, txFunc)
}
