package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"todo/src/controllers"
	"todo/src/entity"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestGetActualNote(t *testing.T) {
	minioServiceImplTest := MinioServiceImplTest{}

	txFunc := func(testJdbcTemplate JdbcTemplateImplTest) {
		var noteService = entity.NoteServiceImpl{
			JdbcTemplate: &testJdbcTemplate,
			MinioService: &minioServiceImplTest}
		
		var note, err = createAndGetNewNote(t, &noteService)
		if err != nil {
			t.Fatalf("ошибка создания note: %s", err)
		}

		var router = createTestRouter(&noteService)

		data := url.Values{}
		data.Set("guid", *note.NoteGuid)
	
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/notes/getActualNote", strings.NewReader(data.Encode()))
		if err != nil {
			t.Fatalf("ошибка формирования http запроса: %s", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	
		router.ServeHTTP(w, req)
	
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

func createTestRouter(noteService *entity.NoteServiceImpl) *gin.Engine {
	var setupTestMiddleware = func (router *gin.Engine) {
		router.Use(note_controller.ApiMiddleware(noteService))
	}
	return note_controller.SetupRouter(setupTestMiddleware)
}

func createAndGetNewNote(t *testing.T, noteService *entity.NoteServiceImpl) (*entity.Note, error) {
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