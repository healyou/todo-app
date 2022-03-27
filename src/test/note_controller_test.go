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
	"todo/src/controllers"
	"todo/src/di"
	"todo/src/entity"
	"todo/src/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRestGetActualNote(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	var router = createTestRouter(di.GetInstance().GetNoteService())

	/* Создаём запрос в rest */
	var note, err = createAndGetNewNote(t, di.GetInstance().GetNoteService())
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
	req.Header.Add(utils.JSON_IN_ACCESS_TOKEN_CODE, CreateTestSuccessTokenWithoutPrivileges(t))

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

func TestRestSaveNoteWithPrivileges(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	/* Создаём запрос в rest */
	res := executeSaveNoteRequest(t, true)
	defer res.Body.Close()
	
	/* Парсим ответ */
	got := ParseResponseBody(t, res)

	/* Проверяем результат */
	var want = gin.H{"result": true}
	assert.Equal(t, want, got)
}

func TestRestSaveNoteWithoutPrivileges(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	/* Выполняем запрос */
	res := executeSaveNoteRequest(t, false)
	defer res.Body.Close()
	
	/* Парсим ответ */
	got := ParseResponseBody(t, res)

	/* Проверяем результат */
	assert.Equal(t, res.StatusCode, http.StatusForbidden)
	_, ok := got["error"]
	assert.True(t, ok, "не найден тег 'error' в json ответе")
}

func executeSaveNoteRequest(t *testing.T, withPrivilege bool) *http.Response {
	var router = createTestRouter(di.GetInstance().GetNoteService())
	
	var note = CreateNewRandomNote()
	noteJsonBytes, err := json.Marshal(*note)
	if err != nil {
		t.Fatalf("ошибка формирования json: %s", err)
	}

	req, err := http.NewRequest("POST", "/notes/saveNote", bytes.NewBuffer(noteJsonBytes))
	if err != nil {
		t.Fatalf("ошибка формирования http запроса: %s", err)
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(noteJsonBytes)))
	req.Header.Set("Content-Type", "application/json")
	if withPrivilege {
		req.Header.Add(utils.JSON_IN_ACCESS_TOKEN_CODE, CreateTestSuccessTokenWithAllPrivileges(t))
	} else {
		req.Header.Add(utils.JSON_IN_ACCESS_TOKEN_CODE, CreateTestSuccessTokenWithoutPrivileges(t))
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	return w.Result()
}

func TestRestDownNoteVersionWithPrivileges(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	/* Выполняем запрос */
	res := executeDownNoteVersionRequest(t, true)
	defer res.Body.Close()
	
	/* Парсим ответ */
	got := ParseResponseBody(t, res)

	/* Проверяем результат */
	var want = gin.H{"result": true}
	assert.Equal(t, want, got)
}

func TestRestDownNoteVersionWithoutPrivileges(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	/* Выполняем запрос */
	res := executeDownNoteVersionRequest(t, false)
	defer res.Body.Close()
	
	/* Парсим ответ */
	got := ParseResponseBody(t, res)

	/* Проверяем результат */
	assert.Equal(t, res.StatusCode, http.StatusForbidden)
	_, ok := got["error"]
	assert.True(t, ok, "не найден тег 'error' в json ответе")
}

func executeDownNoteVersionRequest(t *testing.T, withPrivilege bool) *http.Response {
	noteService := di.GetInstance().GetNoteService()
	var router = createTestRouter(noteService)

	/* Создаём запрос в rest */
	var note = CreateAndGetNewNoteWithNVersion(t, noteService, 2)
	data := url.Values{}
	data.Set("guid", *note.NoteGuid)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/notes/downNoteVersion", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("ошибка формирования http запроса: %s", err)
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if withPrivilege {
		req.Header.Add(utils.JSON_IN_ACCESS_TOKEN_CODE, CreateTestSuccessTokenWithAllPrivileges(t))
	} else {
		req.Header.Add(utils.JSON_IN_ACCESS_TOKEN_CODE, CreateTestSuccessTokenWithoutPrivileges(t))
	}

	/* Выполняем запрос */
	router.ServeHTTP(w, req)

	return w.Result()
}

func TestRestUpNoteVersionWithPrivileges(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	/* Выполняем запрос */
	res := executeUpNoteVersionRequest(t, true)
	defer res.Body.Close()
	
	/* Парсим ответ */
	got := ParseResponseBody(t, res)

	/* Проверяем результат */
	var want = gin.H{"result": true}
	assert.Equal(t, want, got)
}

func TestRestUpNoteVersionWithoutPrivileges(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	/* Выполняем запрос */
	res := executeUpNoteVersionRequest(t, false)
	defer res.Body.Close()
	
	/* Парсим ответ */
	got := ParseResponseBody(t, res)

	/* Проверяем результат */
	assert.Equal(t, res.StatusCode, http.StatusForbidden)
	_, ok := got["error"]
	assert.True(t, ok, "не найден тег 'error' в json ответе")
}

func executeUpNoteVersionRequest(t *testing.T, withPrivilege bool) *http.Response {
	var noteService = di.GetInstance().GetNoteService()
	var router = createTestRouter(noteService)

	/* Создаём запрос в rest */
	var note = CreateAndGetNewNoteWithNVersion(t, noteService, 2)
	err := noteService.DownNoteVersion(*note.NoteGuid)
	if err != nil {
		t.Fatalf("ошибка создания note: %s", err)
	}
	note, err = noteService.GetActualNoteByGuid(*note.NoteGuid)
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
	if withPrivilege {
		req.Header.Add(utils.JSON_IN_ACCESS_TOKEN_CODE, CreateTestSuccessTokenWithAllPrivileges(t))
	} else {
		req.Header.Add(utils.JSON_IN_ACCESS_TOKEN_CODE, CreateTestSuccessTokenWithoutPrivileges(t))
	}

	/* Выполняем запрос */
	router.ServeHTTP(w, req)

	return w.Result()
}

func TestRestGetUserNotes(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	var router = createTestRouter(di.GetInstance().GetNoteService())

	/* Создаём запрос в rest */
	var note, err = createAndGetNewNote(t, di.GetInstance().GetNoteService())
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
	req.Header.Add(utils.JSON_IN_ACCESS_TOKEN_CODE, CreateTestSuccessTokenWithoutPrivileges(t))

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

func createTestRouter(noteService entity.NoteService) *gin.Engine {
	return controllers.SetupRouter()
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

func TestRestGetNoteVersionHistoryWithoutPrivileges(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	/* Выполняем запрос */
	res := executeGetNoteVersionHistoryRequest(t, false, false)
	defer res.Body.Close()
	
	/* Парсим ответ */
	got := ParseResponseBody(t, res)

	/* Проверяем результат */
	assert.Equal(t, res.StatusCode, http.StatusForbidden)
	_, ok := got["error"]
	assert.True(t, ok, "не найден тег 'error' в json ответе")
}

func TestRestGetNoteVersionHistoryNoItems(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	/* Выполняем запрос */
	res := executeGetNoteVersionHistoryRequest(t, true, true)
	defer res.Body.Close()
	
	/* Парсим ответ */
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("ошибка чтения ответа: %s", err)
	}
	var got []entity.NoteVersionInfo
	err = json.Unmarshal(body, &got)
	if err != nil {
		t.Fatalf("ошибка формирования json: %s", err)
	}

	/* Проверяем результат */
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, 0, len(got))
}

func TestRestGetNoteVersionHistoryWithItems(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	/* Выполняем запрос */
	res := executeGetNoteVersionHistoryRequest(t, true, false)
	defer res.Body.Close()
	
	/* Парсим ответ */
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("ошибка чтения ответа: %s", err)
	}
	var got []entity.NoteVersionInfo
	err = json.Unmarshal(body, &got)
	if err != nil {
		t.Fatalf("ошибка формирования json: %s", err)
	}

	/* Проверяем результат */
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.True(t, len(got) > 0)
}

func executeGetNoteVersionHistoryRequest(t *testing.T, withPrivilege bool, randomGuid bool) *http.Response {
	var noteService = di.GetInstance().GetNoteService()
	var router = createTestRouter(noteService)

	/* Создаём запрос в rest */
	var note = CreateAndGetNewNoteWithNVersion(t, noteService, 2)

	data := url.Values{}
	if randomGuid {
		data.Set("guid", uuid.New().String())
	} else {
		data.Set("guid", *note.NoteGuid)
	}
	

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/notes/getNoteVersionHistory", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("ошибка формирования http запроса: %s", err)
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if withPrivilege {
		req.Header.Add(utils.JSON_IN_ACCESS_TOKEN_CODE, CreateTestSuccessTokenWithAllPrivileges(t))
	} else {
		req.Header.Add(utils.JSON_IN_ACCESS_TOKEN_CODE, CreateTestSuccessTokenWithoutPrivileges(t))
	}

	/* Выполняем запрос */
	router.ServeHTTP(w, req)

	return w.Result()
}