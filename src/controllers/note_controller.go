package note_controller

import (
	"log"
	"net/http"
	"strconv"
	"todo/src/di"
	"todo/src/entity"

	"github.com/gin-gonic/gin"
)

func SaveNote(c *gin.Context) {
	var note entity.Note
	var err = c.BindJSON(&note)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

	var noteService = di.GetInstance().GetNoteService()

	_, err = noteService.SaveNote(&note)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"result": true})
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
		log.Println(err)
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	var noteService = di.GetInstance().GetNoteService()

	notes, err := noteService.GetUserActualNotes(userIdParam)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, notes)
}

func DownNoteVersion(c *gin.Context) {
	var guidParam string = c.PostForm("guid")
	if len(guidParam) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Не указан параметр 'guid'"})
		return
	}

	var noteService = di.GetInstance().GetNoteService()

	var err = noteService.DownNoteVersion(guidParam)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"result": true})
}

func UpNoteVersion(c *gin.Context) {
	var guidParam string = c.PostForm("guid")
	if len(guidParam) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Не указан параметр 'guid'"})
		return
	}

	var noteService = di.GetInstance().GetNoteService()

	var err = noteService.UpNoteVersion(guidParam)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"result": true})
}
