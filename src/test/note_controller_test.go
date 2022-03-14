package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	note_controller "todo/src/controllers"
	"todo/src/di"
	"todo/src/entity"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestRestGetActualNote(t *testing.T) {
	txFunc := func(di di.DependencyInjection) {
		var router = createTestRouter(di.GetNoteService())

		/* Создаём запрос в rest */
		var note, err = createAndGetNewNote(t, di.GetNoteService())
		if err != nil {
			t.Fatalf("ошибка создания note: %s", err)
		}
		data := url.Values{}
		data.Set("guid", *note.NoteGuid)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/notes/getActualNote", strings.NewReader(data.Encode()))
		if err != nil {
			t.Fatalf("ошибка формирования http запроса: %s", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

		/* Выполняем запрос */
		router.ServeHTTP(w, req)

		/* Проверяем результат */
		var want gin.H
		wantBytes, err := json.Marshal(*note)
		if err != nil {
			t.Fatalf("ошибка формирования json: %s", err)
		}
		json.Unmarshal(wantBytes, &want)

		var got gin.H
		err = json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatalf("ошибка формирования json: %s", err)
		}
		assert.Equal(t, want, got)
	}

	ExecuteTestRollbackTransaction(t, txFunc)
}

func TestRestSaveNote(t *testing.T) {
	txFunc := func(di di.DependencyInjection) {
		var router = createTestRouter(di.GetNoteService())

		/* Создаём запрос в rest */
		var note = CreateNewRandomNote()
		noteJsonBytes, err := json.Marshal(*note)
		if err != nil {
			t.Fatalf("ошибка формирования json: %s", err)
		}

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/notes/saveNote", bytes.NewBuffer(noteJsonBytes))
		if err != nil {
			t.Fatalf("ошибка формирования http запроса: %s", err)
		}
		req.Header.Add("Content-Length", strconv.Itoa(len(noteJsonBytes)))
		req.Header.Set("Content-Type", "application/json")

		/* Выполняем запрос */
		router.ServeHTTP(w, req)

		/* Проверяем результат */
		var want = gin.H{"result": true}
		var got gin.H
		err = json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatalf("ошибка формирования json: %s", err)
		}
		assert.Equal(t, want, got)
	}

	ExecuteTestRollbackTransaction(t, txFunc)
}

func TestRestDownNoteVersion(t *testing.T) {
	txFunc := func(di di.DependencyInjection) {
		var router = createTestRouter(di.GetNoteService())

		/* Создаём запрос в rest */
		var note, err = createAndGetNewNoteWith2Version(t, di.GetNoteService())
		if err != nil {
			t.Fatalf("ошибка создания note: %s", err)
		}
		data := url.Values{}
		data.Set("guid", *note.NoteGuid)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/notes/downNoteVersion", strings.NewReader(data.Encode()))
		if err != nil {
			t.Fatalf("ошибка формирования http запроса: %s", err)
		}
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		/* Выполняем запрос */
		router.ServeHTTP(w, req)

		/* Проверяем результат */
		var want = gin.H{"result": true}
		var got gin.H
		err = json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatalf("ошибка формирования json: %s", err)
		}
		assert.Equal(t, want, got)
	}

	ExecuteTestRollbackTransaction(t, txFunc)
}

func TestRestUpNoteVersion(t *testing.T) {
	txFunc := func(di di.DependencyInjection) {
		var router = createTestRouter(di.GetNoteService())

		/* Создаём запрос в rest */
		var note, err = createAndGetNewNoteWith2Version(t, di.GetNoteService())
		if err != nil {
			t.Fatalf("ошибка создания note: %s", err)
		}
		err = di.GetNoteService().DownNoteVersion(*note.NoteGuid)
		if err != nil {
			t.Fatalf("ошибка создания note: %s", err)
		}
		note, err = di.GetNoteService().GetActualNoteByGuid(*note.NoteGuid)
		if err != nil {
			t.Fatalf("ошибка создания note: %s", err)
		}

		data := url.Values{}
		data.Set("guid", *note.NoteGuid)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/notes/upNoteVersion", strings.NewReader(data.Encode()))
		if err != nil {
			t.Fatalf("ошибка формирования http запроса: %s", err)
		}
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		/* Выполняем запрос */
		router.ServeHTTP(w, req)

		/* Проверяем результат */
		var want = gin.H{"result": true}
		var got gin.H
		err = json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatalf("ошибка формирования json: %s", err)
		}
		assert.Equal(t, want, got)
	}

	ExecuteTestRollbackTransaction(t, txFunc)
}

func TestRestGetUserNotes(t *testing.T) {
	txFunc := func(di di.DependencyInjection) {
		var router = createTestRouter(di.GetNoteService())

		/* Создаём запрос в rest */
		var note, err = createAndGetNewNote(t, di.GetNoteService())
		if err != nil {
			t.Fatalf("ошибка создания note: %s", err)
		}

		data := url.Values{}
		data.Set("user_id", strconv.FormatInt(*note.UserId, 10))

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/notes/getUserNotes", strings.NewReader(data.Encode()))
		if err != nil {
			t.Fatalf("ошибка формирования http запроса: %s", err)
		}
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		/* Выполняем запрос */
		router.ServeHTTP(w, req)

		/* Проверяем результат */
		body, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Fatalf("ошибка чтения ответа: %s", err)
		}
		var got []entity.Note
		err = json.Unmarshal(body, &got)
		if err != nil {
			t.Fatalf("ошибка формирования json: %s", err)
		}

		assert.True(t, len(got) >= 1)
		var findedRecord = false
		for i := 0; i < len(got); i++ {
			gotNote := got[i]
			if (*gotNote.NoteGuid == *note.NoteGuid) {
				findedRecord = true
				assert.Equal(t, *note.Id, *gotNote.Id)
			}
		}
		assert.Equal(t, findedRecord, true)
	}

	ExecuteTestRollbackTransaction(t, txFunc)
}

func createTestRouter(noteService entity.NoteService) *gin.Engine {
	return note_controller.SetupRouter()
}

func createAndGetNewNote(t *testing.T, noteService entity.NoteService) (*entity.Note, error) {
	noteId, err := noteService.SaveNote(CreateNewRandomNote())
	if err != nil {
		t.Fatalf("error was not expected while test method: %s", err)
	}
	note, err := noteService.GetNote(*noteId)
	if err != nil {
		t.Fatalf("error was not expected while test method: %s", err)
	}
	return note, nil
}

func createAndGetNewNoteWith2Version(t *testing.T, noteService entity.NoteService) (*entity.Note, error) {
	/* Создаём новый note */
	noteId, err := noteService.SaveNote(CreateNewRandomNote())
	if err != nil {
		t.Fatalf("error was not expected while test method: %s", err)
	}
	note, err := noteService.GetNote(*noteId)
	if err != nil {
		t.Fatalf("error was not expected while test method: %s", err)
	}
	/* Создаём 2 версию note */
	noteId, err = noteService.SaveNote(note)
	if err != nil {
		t.Fatalf("error was not expected while test method: %s", err)
	}
	note, err = noteService.GetNote(*noteId)
	if err != nil {
		t.Fatalf("error was not expected while test method: %s", err)
	}
	return note, nil
}