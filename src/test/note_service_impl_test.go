package test

import (
	"testing"
	"todo/src/di"
	"todo/src/entity"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSaveNewNoteWithData(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	savedNote := CreateNewRandomNote()

	var noteService = di.GetInstance().GetNoteService()
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
	assert.Nil(t, createdNote.PrevNoteVersionId)
	assert.Equal(t, *createdNote.Title, *savedNote.Title)
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

func TestGetNote(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	var noteService = di.GetInstance().GetNoteService()

	noteId, err := noteService.SaveNote(CreateNewRandomNote())
	if err != nil {
		t.Fatalf("error was not expected while test method: %s", err)
	}
	
	result, err := noteService.GetNote(*noteId)
	if err != nil {
		t.Fatalf("error was not expected while test method: %s", err)
	}

	assert.Equal(t, *result.Id, *noteId)
	assert.NotNil(t, *result.NoteGuid)
	assert.NotNil(t, *result.Version)
	assert.NotNil(t, *result.Title)
	assert.NotNil(t, *result.Text)
	assert.Nil(t, result.PrevNoteVersionId)
	assert.NotNil(t, *result.UserId)
	assert.NotNil(t, *result.CreateDate)
	assert.NotNil(t, *result.Deleted)
	assert.NotNil(t, *result.Archive)
	assert.NotNil(t, result.NoteFiles)
}

func TestGetNoteByNoteGuid(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	var noteService = di.GetInstance().GetNoteService()

	var note = CreateNewRandomNote()
	_, err := noteService.SaveNote(note)
	if err != nil {
		t.Fatalf("error was not expected while test method: %s", err)
	}

	result, err := noteService.GetActualNoteByGuid(*note.NoteGuid)
	if err != nil {
		t.Fatalf("error was not expected while test method: %s", err)
	}

	assert.NotNil(t, *result.Id)
	assert.Equal(t, *result.NoteGuid, *result.NoteGuid)
	assert.Equal(t, *result.Version, *result.Version)
	assert.Equal(t, *result.Title, *result.Title)
	assert.Equal(t, *result.Text, *result.Text)
	assert.Equal(t, *result.UserId, *result.UserId)
	assert.Nil(t, result.PrevNoteVersionId)
	assert.NotNil(t, *result.CreateDate)
	assert.Equal(t, *result.Deleted, *result.Deleted)
	assert.Equal(t, *result.Archive, *result.Archive)
	assert.NotNil(t, result.NoteFiles)
}

func TestUpdateNote(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	var noteService = di.GetInstance().GetNoteService()

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
	assert.Equal(t, *note.Version+1, *updatedNote.Version)
	assert.Equal(t, *note.Title, *updatedNote.Title)
	assert.Equal(t, *note.Text, *updatedNote.Text)
	assert.Equal(t, *note.UserId, *updatedNote.UserId)
	assert.Equal(t, *note.Deleted, *updatedNote.Deleted)
	assert.Equal(t, *note.Archive, *updatedNote.Archive)
	assert.Equal(t, *updatedNote.Actual, true)
	assert.NotNil(t, updatedNote.NoteFiles)
	assert.Equal(t, *note.Id, *updatedNote.PrevNoteVersionId)
	assert.Equal(t, 2, len(updatedNote.NoteFiles))
}

func TestDownNoteVersion(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	var noteService = di.GetInstance().GetNoteService()

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
	assert.Equal(t, *firstNoteVersion.Title, *secondNoteVersion.Title)
	assert.NotEqual(t, *firstNoteVersion.Text, *secondNoteVersion.Text)
	assert.Equal(t, *firstNoteVersionId, *secondNoteVersion.PrevNoteVersionId)
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
	assert.Nil(t, downGradeNote.PrevNoteVersionId)
	assert.Equal(t, downGradeNote.PrevNoteVersionId, firstNoteVersion.PrevNoteVersionId)
	assert.Equal(t, *downGradeNote.NoteGuid, *note.NoteGuid)
	assert.Equal(t, *downGradeNote.Id, *firstNoteVersion.Id)
	assert.Equal(t, *downGradeNote.Title, *firstNoteVersion.Title)
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

func TestErrorDownNewNoteVersion(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	var noteService = di.GetInstance().GetNoteService()

	/* Создаём 1 версию note */
	note := CreateNewRandomNote()
	_, err := noteService.SaveNote(note)
	if err != nil {
		t.Fatalf("не удалось сохранить note: %s", err)
	}
	firstNoteVersion, err := noteService.GetActualNoteByGuid(*note.NoteGuid)
	if err != nil {
		t.Fatalf("не удалось получить note: %s", err)
	}

	/* Пытаемся уменьшить версию */
	err = noteService.DownNoteVersion(*firstNoteVersion.NoteGuid)
	if err == nil {
		t.Fatalf("получилось уменьшить версию, а такого быть не должно")
	}
	assert.NotNil(t, err)
}

func TestErrorDoubleDownNewNoteVersion(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	var noteService = di.GetInstance().GetNoteService()

	/* Создаём 1 версию note */
	note := CreateNewRandomNote()
	_, err := noteService.SaveNote(note)
	if err != nil {
		t.Fatalf("не удалось сохранить note: %s", err)
	}
	firstNoteVersion, err := noteService.GetActualNoteByGuid(*note.NoteGuid)
	if err != nil {
		t.Fatalf("не удалось получить note: %s", err)
	}

	/* Создаём вторую версию note */
	_, err = noteService.SaveNote(firstNoteVersion)
	if err != nil {
		t.Fatalf("не удалось обновить note: %s", err)
	}

	/* Уменьшаем версию note */
	err = noteService.DownNoteVersion(*firstNoteVersion.NoteGuid)
	if err != nil {
		t.Fatalf("не удалось уменьшить версию note: %s", err)
	}

	/* Второй раз должна быть ошибка, т.к. некуда уменьшать версию */
	err = noteService.DownNoteVersion(*firstNoteVersion.NoteGuid)
	if err == nil {
		t.Fatalf("получилось уменьшить версию, а такого быть не должно")
	}
	assert.NotNil(t, err)
}

func TestUpNoteVersion(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	var noteService = di.GetInstance().GetNoteService()

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
	_, err = noteService.SaveNote(firstNoteVersion)
	if err != nil {
		t.Fatalf("не удалось сохранить note: %s", err)
	}
	secondNoteVersion, err := noteService.GetActualNoteByGuid(*note.NoteGuid)
	if err != nil {
		t.Fatalf("не удалось получить note: %s", err)
	}
	_, err = noteService.GetNote(*firstNoteVersionId)
	if err != nil {
		t.Fatalf("не удалось получить note: %s", err)
	}

	/* уменьшаем версию note */
	err = noteService.DownNoteVersion(*secondNoteVersion.NoteGuid)
	if err != nil {
		t.Fatalf("не удалось уменьшить версию note: %s", err)
	}
	downGradeNote, err := noteService.GetActualNoteByGuid(*secondNoteVersion.NoteGuid)
	if err != nil {
		t.Fatalf("не удалось получить note: %s", err)
	}

	/* Поднимаем версию note обратно */
	err = noteService.UpNoteVersion(*downGradeNote.NoteGuid)
	if err != nil {
		t.Fatalf("не удалось увеличить версию note: %s", err)
	}
	upGradeNote, err := noteService.GetActualNoteByGuid(*downGradeNote.NoteGuid)
	if err != nil {
		t.Fatalf("не удалось получить note: %s", err)
	}

	/* Проверяем, что значения полей стали снова равны предыдущей версии */
	/* Проверяем актуальность данных c первой версией */
	assert.Equal(t, *upGradeNote.PrevNoteVersionId, *secondNoteVersion.PrevNoteVersionId)
	assert.Equal(t, *upGradeNote.NoteGuid, *secondNoteVersion.NoteGuid)
	assert.Equal(t, *upGradeNote.Id, *secondNoteVersion.Id)
	assert.Equal(t, *upGradeNote.Title, *secondNoteVersion.Title)
	assert.Equal(t, *upGradeNote.Text, *secondNoteVersion.Text)
	assert.Equal(t, *upGradeNote.Version, *secondNoteVersion.Version)
	assert.Equal(t, *upGradeNote.Archive, *secondNoteVersion.Archive)
	assert.Equal(t, *upGradeNote.Deleted, *secondNoteVersion.Deleted)
	assert.Equal(t, *upGradeNote.CreateDate, *secondNoteVersion.CreateDate)
	assert.Equal(t, *upGradeNote.UserId, *secondNoteVersion.UserId)
	assert.Equal(t, firstNoteFilesCount, len(secondNoteVersion.NoteFiles))
	assert.Equal(t, firstNoteFilesCount, len(upGradeNote.NoteFiles))
	assert.True(t, *downGradeNote.Actual)
}

func TestErrorUpNewNoteVersion(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	var noteService = di.GetInstance().GetNoteService()

	/* Создаём 1 версию note */
	note := CreateNewRandomNote()
	_, err := noteService.SaveNote(note)
	if err != nil {
		t.Fatalf("не удалось сохранить note: %s", err)
	}
	firstNoteVersion, err := noteService.GetActualNoteByGuid(*note.NoteGuid)
	if err != nil {
		t.Fatalf("не удалось получить note: %s", err)
	}

	/* Пытаемся увеличить версию */
	err = noteService.UpNoteVersion(*firstNoteVersion.NoteGuid)
	if err == nil {
		t.Fatalf("получилось увеличить версию, а такого быть не должно")
	}
	assert.NotNil(t, err)
}

func TestErrorDoubleUpNewNoteVersion(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	var noteService = di.GetInstance().GetNoteService()

	/* Создаём 1 версию note */
	note := CreateNewRandomNote()
	_, err := noteService.SaveNote(note)
	if err != nil {
		t.Fatalf("не удалось сохранить note: %s", err)
	}
	firstNoteVersion, err := noteService.GetActualNoteByGuid(*note.NoteGuid)
	if err != nil {
		t.Fatalf("не удалось получить note: %s", err)
	}

	/* Создаём вторую версию note */
	_, err = noteService.SaveNote(firstNoteVersion)
	if err != nil {
		t.Fatalf("не удалось обновить note: %s", err)
	}

	/* Уменьшаем версию note */
	err = noteService.DownNoteVersion(*firstNoteVersion.NoteGuid)
	if err != nil {
		t.Fatalf("не удалось уменьшить версию note: %s", err)
	}

	/* Увеличиваем версию note */
	err = noteService.UpNoteVersion(*firstNoteVersion.NoteGuid)
	if err != nil {
		t.Fatalf("не удалось уменьшить версию note: %s", err)
	}

	/* Пытаемся увеличить версию, но должна быть ошибка, т.к. некуда увеличивать версию */
	err = noteService.UpNoteVersion(*firstNoteVersion.NoteGuid)
	if err == nil {
		t.Fatalf("получилось увеличить версию, а такого быть не должно")
	}
	assert.NotNil(t, err)
}


func TestGetUserNotes(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	var noteService = di.GetInstance().GetNoteService()

	/* Создаём notes */
	notes := [2]entity.Note{*CreateNewRandomNote(), *CreateNewRandomNote()}
	userId := notes[0].UserId

	for index, note := range notes {
		noteId, err := noteService.SaveNote(&note)
		note.Id = noteId
		if err != nil {
			t.Fatalf("не удалось сохранить note: %s", err)
		}
		notes[index] = note
	}

	/* Получаем записи */
	userNotes, err := noteService.GetUserActualNotes(*userId)
	if err != nil {
		t.Fatalf("не получилось получить notes: %s", err)
	}

	/* Проверяем, что данные начитались */
	assert.NotNil(t, userNotes)
	assert.True(t, len(userNotes) > 0)
	for _, savedNote := range userNotes {
		for _, createNote := range notes {
			if (*savedNote.Id == *createNote.Id) {
				assert.Equal(t, len(createNote.NoteFiles), len(savedNote.NoteFiles))
				assert.Equal(t, *userId, *savedNote.UserId)
				assert.Equal(t, *createNote.UserId, *savedNote.UserId)
			}
		}

		assert.Equal(t, *savedNote.Actual, true)
	}
}

func TestGetNoteVersionHistoryNoItems(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	var noteService = di.GetInstance().GetNoteService()
	history, err := noteService.GetNoteVersionHistory(uuid.New().String())

	assert.Nil(t, err)
	assert.Equal(t, len(history), 0)
}

func TestGetNoteVersionHistoryWithItems(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	var countHistoryItems = 10
	var noteService = di.GetInstance().GetNoteService()
	note := CreateAndGetNewNoteWithNVersion(t, noteService, countHistoryItems)
	history, err := noteService.GetNoteVersionHistory(*note.NoteGuid)
	if err != nil {
		t.Fatalf("не получилось получить истории версий note: %s", err)
	}

	assert.Nil(t, err)
	assert.Equal(t, countHistoryItems, len(history))
}