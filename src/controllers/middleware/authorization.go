package middleware

import (
	"fmt"
	"log"
	"net/http"
	"todo/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

func AuthorizationMiddleware(c *gin.Context) {
	userAuthData, err := parseUserAuthDataFromAccessToken(c)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
	}

	c.Set(utils.GIN_CONTEXT_USER_AUTH_DATA, userAuthData)
	c.Next()
}

func parseUserAuthDataFromAccessToken(c *gin.Context) (*UserAuthData, error) {
	var accessTokenHeader []string = c.Request.Header[utils.JSON_IN_ACCESS_TOKEN_CODE]
	if len(accessTokenHeader) == 0 || len(accessTokenHeader[0]) == 0 {
		return nil, errors.New("не найден токен пользователя")
	}

	userAccessToken := accessTokenHeader[0]
	token, _, err := new(jwt.Parser).ParseUnverified(userAccessToken, &UserAuthData{})
	if err != nil {
		return nil, errors.Wrap(err, "ошибка парсинга токена")
	}

	userAuthData, ok := token.Claims.(*UserAuthData)
	if !ok {
		return nil, errors.New("token.Claims is not cast to '*UserAuthData'")
	}

	err = userAuthData.Valid()
	if err != nil {
		return nil, err
	}

	return userAuthData, nil
}