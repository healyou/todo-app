package test_utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

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