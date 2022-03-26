package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
	"todo/src/di"
	"todo/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestNoAccessTokenError(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	/* Выполняем тестовый запрос */
	res := executeTestGetActualNoteRequest(t, func(header *http.Header) {
		/* no token */
	})
	defer res.Body.Close()

	/* Парсим ответ */
	got := parseResponseBody(t, res)

	/* Проверяем результат */
	assertErrorTokenStatus(t, res, got)
}

func TestNotValidAccessTokenError(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	/* Выполняем тестовый запрос */
	res := executeTestGetActualNoteRequest(t, func(header *http.Header) {
		/* not valid access token */
		header.Add(utils.JSON_IN_ACCESS_TOKEN_CODE, "not token value")
	})
	defer res.Body.Close()

	/* Парсим ответ */
	got := parseResponseBody(t, res)

	/* Проверяем результат */
	assertErrorTokenStatus(t, res, got)
}

func TestNotValidTokenClaimsError(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	/* Выполняем тестовый запрос */
	res := executeTestGetActualNoteRequest(t, func(header *http.Header) {
		/* not valid access token */
		header.Add(utils.JSON_IN_ACCESS_TOKEN_CODE, createTestTokenWithoutUsername(t))
	})
	defer res.Body.Close()

	/* Парсим ответ */
	got := parseResponseBody(t, res)

	/* Проверяем результат */
	assertErrorTokenStatus(t, res, got)
}

func TestValidTokenSuccess(t *testing.T) {
	closeIntegrationTest := InitIntegrationTest(t)
	defer closeIntegrationTest(t)

	/* Выполняем тестовый запрос */
	res := executeTestGetActualNoteRequest(t, func(header *http.Header) {
		/* valid token */
		header.Add(utils.JSON_IN_ACCESS_TOKEN_CODE, CreateTestSuccessToken(t))
	})
	defer res.Body.Close()

	/* Парсим ответ */
	got := parseResponseBody(t, res)

	/* Проверяем результат */
	assert.Equal(t, res.StatusCode, http.StatusOK)
	_, ok := got["error"]
	assert.False(t, ok, "не найден тег 'error' в json ответе")
}

func executeTestGetActualNoteRequest(t *testing.T, headerModifier func(header *http.Header)) *http.Response {
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
	headerModifier(&req.Header)

	/* Выполняем запрос */
	router.ServeHTTP(w, req)
	/* Получаем результат */
	return w.Result()
}

func parseResponseBody(t *testing.T, res *http.Response) gin.H {
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

func CreateTestSuccessToken(t *testing.T) string {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["privileges"] = []string{}
	atClaims["user_id"] = 1
	atClaims["username"] = "admin"
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte("jdnfksdmfksd"))
	if err != nil {
		t.Fatalf("ошибка формирования токена: %s", err)
	}
	return token
}

func createTestTokenWithoutUsername(t *testing.T) string {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["privileges"] = []string{}
	atClaims["user_id"] = 1
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte("jdnfksdmfksd"))
	if err != nil {
		t.Fatalf("ошибка формирования токена: %s", err)
	}
	return token
}

func assertErrorTokenStatus(t *testing.T, res *http.Response, got gin.H) {
	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
	_, ok := got["error"]
	assert.True(t, ok, "не найден тег 'error' в json ответе")
}