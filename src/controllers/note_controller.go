package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"todo/src/controllers/middleware"
	"todo/src/di"
	"todo/src/entity"
	"todo/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func DownloadNoteFile(c *gin.Context) {
	_, err := getUserAuthData(c)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var idParamStr string = c.PostForm("id")
	if len(idParamStr) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Не указан параметр 'id'"})
		return
	}

	idParam, err := strconv.ParseInt(idParamStr, 10, 64)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	var noteService = di.GetInstance().GetNoteService()
	noteFile, err := noteService.GetNoteFile(idParam)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.Header("Content-Disposition", "attachment; filename=\""+*noteFile.Filename+"\"")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Length", strconv.Itoa(len(noteFile.Data)))
	c.Data(http.StatusOK, "application/octet-stream", noteFile.Data)
}

func SaveNote(c *gin.Context) {
	userAuthData, err := getUserAuthData(c)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !userAuthData.HasPrivilege(middleware.CREATE_NOTE_PRIVILEGE) {
		message := "недостаточно привилегий для выполнения операции"
		log.Println("у пользователя - " + *userAuthData.Username + " " + message)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": message})
		return
	}

	var note entity.Note
	err = c.BindJSON(&note)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

	var noteService = di.GetInstance().GetNoteService()

	_, err = noteService.SaveNote(&note)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"result": true})
}

func getUserAuthData(c *gin.Context) (*middleware.UserAuthData, error) {
	data, ok := c.Get(utils.GIN_CONTEXT_USER_AUTH_DATA)
	if !ok {
		return nil, errors.New("не найдены данные авторизации пользователя " + utils.GIN_CONTEXT_USER_AUTH_DATA)
	}
	userAuthData, ok := data.(*middleware.UserAuthData)
	if !ok {
		return nil, errors.New("данные авторизации пользователя не являются типом 'UserAuthData'")
	}

	return userAuthData, nil
}

func GetActualNote(c *gin.Context) {
	var guidParam string = c.PostForm("guid")
	if len(guidParam) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Не указан параметр 'guid'"})
		return
	}

	var noteService = di.GetInstance().GetNoteService()

	note, err := noteService.GetActualNoteByGuid(guidParam)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, note)
}

func GetUserNotes(c *gin.Context) {
	var userIdParamStr string = c.PostForm("user_id")
	if len(userIdParamStr) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Не указан параметр 'user_id'"})
		return
	}

	var userIdParam, err = strconv.ParseInt(userIdParamStr, 10, 64)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	var noteService = di.GetInstance().GetNoteService()

	notes, err := noteService.GetUserActualNotes(userIdParam)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, notes)
}

func GetLastUserNoteMainInfo(c *gin.Context) {
	var userIdParamStr string = c.PostForm("user_id")
	if len(userIdParamStr) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Не указан параметр 'user_id'"})
		return
	}

	var maxCountLimitParamStr string = c.PostForm("max_count_limit")
	if len(maxCountLimitParamStr) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Не указан параметр 'max_count_limit'"})
		return
	}

	var userIdParam, err = strconv.ParseInt(userIdParamStr, 10, 64)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	maxCountLimitParam, err := strconv.ParseInt(maxCountLimitParamStr, 10, 64)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	var noteService = di.GetInstance().GetNoteService()

	notesMainInfo, err := noteService.GetLastUserNoteMainInfo(userIdParam, maxCountLimitParam)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, notesMainInfo)
}

func DownNoteVersion(c *gin.Context) {
	userAuthData, err := getUserAuthData(c)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !userAuthData.HasPrivilege(middleware.CHANGE_NOTE_VERSION_PRIVILEGE) {
		message := "недостаточно привилегий для выполнения операции"
		log.Println("у пользователя - " + *userAuthData.Username + " " + message)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": message})
		return
	}

	var guidParam string = c.PostForm("guid")
	if len(guidParam) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Не указан параметр 'guid'"})
		return
	}

	var noteService = di.GetInstance().GetNoteService()

	err = noteService.DownNoteVersion(guidParam)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"result": true})
}

func UpNoteVersion(c *gin.Context) {
	userAuthData, err := getUserAuthData(c)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !userAuthData.HasPrivilege(middleware.CHANGE_NOTE_VERSION_PRIVILEGE) {
		message := "недостаточно привилегий для выполнения операции"
		log.Println("у пользователя - " + *userAuthData.Username + " " + message)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": message})
		return
	}

	var guidParam string = c.PostForm("guid")
	if len(guidParam) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Не указан параметр 'guid'"})
		return
	}

	var noteService = di.GetInstance().GetNoteService()

	err = noteService.UpNoteVersion(guidParam)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"result": true})
}

func GetNoteVersionHistory(c *gin.Context) {
	userAuthData, err := getUserAuthData(c)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !userAuthData.HasPrivilege(middleware.VIEW_NOTE_VERSION_HISTORY_PRIVILEGE) {
		message := "недостаточно привилегий для выполнения операции"
		log.Println("у пользователя - " + *userAuthData.Username + " " + message)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": message})
		return
	}

	var guidParam string = c.PostForm("guid")
	if len(guidParam) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Не указан параметр 'guid'"})
		return
	}

	var noteService = di.GetInstance().GetNoteService()
	history, err := noteService.GetNoteVersionHistory(guidParam)
	if err != nil {
		log.Println(fmt.Printf("%+v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, history)
}
