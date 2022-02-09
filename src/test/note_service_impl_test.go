package test

import (
	"testing"
	"todo/src/entity"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestSaveNewNoteWithData(t *testing.T) {
	minioServiceImplTest := MinioServiceImplTest{}

	txFunc := func(testJdbcTemplate JdbcTemplateImplTest) {
		savedNote := CreateNewRandomNote()

		var noteService = entity.NoteServiceImpl{
			JdbcTemplate: &testJdbcTemplate,
			MinioService: &minioServiceImplTest}
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
			assert.NotNil(t, createdFile.Data)
			assert.Equal(t, *createdFile.NoteId, *createdNote.Id)
			assert.Equal(t, *createdFile.Filename, *expectedFile.Filename)
		}
	}
	ExecuteTestRollbackTransaction(t, txFunc)
}

func TestGetNote(t *testing.T) {
	minioServiceImplTest := MinioServiceImplTest{}

	txFunc := func(testJdbcTemplate JdbcTemplateImplTest) {
		var noteService = entity.NoteServiceImpl{
			JdbcTemplate: &testJdbcTemplate,
			MinioService: &minioServiceImplTest}
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
	minioServiceImplTest := MinioServiceImplTest{}

	txFunc := func(testJdbcTemplate JdbcTemplateImplTest) {
		expectedNote := GetNoteWithMaxVersion()

		var noteService = entity.NoteServiceImpl{
			JdbcTemplate: &testJdbcTemplate,
			MinioService: &minioServiceImplTest}
		result, err := noteService.GetActualNoteByGuid(*expectedNote.NoteGuid)
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

func TestUpdateNote(t *testing.T) {
	minioServiceImplTest := MinioServiceImplTest{}

	txFunc := func(testJdbcTemplate JdbcTemplateImplTest) {
		var noteService = entity.NoteServiceImpl{
			JdbcTemplate: &testJdbcTemplate,
			MinioService: &minioServiceImplTest}
		
		noteId, err := noteService.SaveNote(CreateNewRandomNote())
		if err != nil {
			t.Fatalf("error was not expected while test method: %s", err)
		}
		note, err := noteService.GetNote(*noteId)
		if err != nil {
			t.Fatalf("error was not expected while test method: %s", err)
		}

		/* Обновляем запись */
		*note.Text = "updatedText"
		updatedfile := note.NoteFiles[1]
		updatedfile.Data = []byte{}
		note.NoteFiles = []entity.NoteFile{
			*createNewRandomNoteFile(),
			updatedfile,
		}
		savedNoteId, err := noteService.SaveNote(note)
		if err != nil {
			t.Fatalf("error was not expected while test method: %s", err)
		}

		/* Смотрим, что сохранилось */
		updatedNote, err := noteService.GetActualNoteByGuid(*note.NoteGuid)
		if err != nil {
			t.Fatalf("error was not expected while test method: %s", err)
		}

		assert.Equal(t, *savedNoteId, *updatedNote.Id)
		assert.Equal(t, *note.NoteGuid, *updatedNote.NoteGuid)
		assert.Equal(t, *note.Version + 1, *updatedNote.Version)
		assert.Equal(t, *note.Text, *updatedNote.Text)
		assert.Equal(t, *note.UserId, *updatedNote.UserId)
		assert.Equal(t, *note.Deleted, *updatedNote.Deleted)
		assert.Equal(t, *note.Archive, *updatedNote.Archive)
		assert.Equal(t, *updatedNote.Actual, true)
		assert.NotNil(t, updatedNote.NoteFiles)
		assert.Equal(t, 2, len(updatedNote.NoteFiles))
	}
	ExecuteTestRollbackTransaction(t, txFunc)
}

func TestDownNoteVersion(t *testing.T) {
	minioServiceImplTest := MinioServiceImplTest{}

	txFunc := func(testJdbcTemplate JdbcTemplateImplTest) {
		var noteService = entity.NoteServiceImpl{
			JdbcTemplate: &testJdbcTemplate,
			MinioService: &minioServiceImplTest}

		/* Создаём 1 версию note */
		note := CreateNewRandomNote()
		firstNoteVersionId, err := noteService.SaveNote(note)
		if err != nil {
			t.Fatalf("не удалось сохранить note: %s", err)
		}
		firstNoteVersion, err := noteService.GetActualNoteByGuid(*note.NoteGuid)
		if err != nil {
			t.Fatalf("не удалось получить note: %s", err)
		}
		firstNoteFilesCount := len(firstNoteVersion.NoteFiles)
		
		/* Создаём 2 версию note */
		*firstNoteVersion.Text = "NOTE 2 VERSION TEXT"
		secondNoteVersionId, err := noteService.SaveNote(firstNoteVersion)
		if err != nil {
			t.Fatalf("не удалось сохранить note: %s", err)
		}
		secondNoteVersion, err := noteService.GetActualNoteByGuid(*note.NoteGuid)
		if err != nil {
			t.Fatalf("не удалось получить note: %s", err)
		}
		firstNoteVersion, err = noteService.GetNote(*firstNoteVersionId)
		if err != nil {
			t.Fatalf("не удалось получить note: %s", err)
		}

		assert.NotEqual(t, *firstNoteVersionId, secondNoteVersionId)
		assert.NotEqual(t, *firstNoteVersion.Text, *secondNoteVersion.Text)
		assert.Equal(t, *secondNoteVersionId, *secondNoteVersion.Id)
		assert.Equal(t, *firstNoteVersion.NoteGuid, *secondNoteVersion.NoteGuid)
		assert.Equal(t, *note.NoteGuid, *secondNoteVersion.NoteGuid)

		/* уменьшаем версию note */
		err = noteService.DownNoteVersion(*secondNoteVersion.NoteGuid)
		if err != nil {
			t.Fatalf("не удалось уменьшить версию note: %s", err)
		}
		downGradeNote, err := noteService.GetActualNoteByGuid(*secondNoteVersion.NoteGuid)
		if err != nil {
			t.Fatalf("не удалось получить note: %s", err)
		}

		/* Проверяем актуальность данных c первой версией */
		assert.Equal(t, *downGradeNote.NoteGuid, *note.NoteGuid)
		assert.Equal(t, *downGradeNote.Id, *firstNoteVersion.Id)
		assert.Equal(t, *downGradeNote.Text, *firstNoteVersion.Text)
		assert.Equal(t, *downGradeNote.Version, *firstNoteVersion.Version)
		assert.Equal(t, *downGradeNote.Archive, *firstNoteVersion.Archive)
		assert.Equal(t, *downGradeNote.Deleted, *firstNoteVersion.Deleted)
		assert.Equal(t, *downGradeNote.CreateDate, *firstNoteVersion.CreateDate)
		assert.Equal(t, *downGradeNote.NoteGuid, *firstNoteVersion.NoteGuid)
		assert.Equal(t, *downGradeNote.UserId, *firstNoteVersion.UserId)
		assert.Equal(t, firstNoteFilesCount, len(downGradeNote.NoteFiles))
		assert.True(t, *downGradeNote.Actual)
		/* Проверяем, что значения не равны второй версии note */
		assert.NotEqual(t, *downGradeNote.Text, *secondNoteVersion.Text)
		assert.NotEqual(t, *downGradeNote.Version, *secondNoteVersion.Version)
	}

	ExecuteTestRollbackTransaction(t, txFunc)
}
