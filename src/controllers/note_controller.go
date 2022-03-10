package note_controller

import (
	"log"
	"net/http"
	"strconv"
	"todo/src/db"
	"todo/src/entity"
	"todo/src/filestorage"
	"todo/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func SaveNote(c *gin.Context) {
	var notes = []entity.Note{
		//{Id: 1, NoteGuid: "not guid", Version: 1,
		//	Text: "text", UserId: 1, Deleted: false, Archive: false,
		//	NoteFiles: []entity.NoteFile{
		//		{Id: 1, NoteId: 1, Guid: "note file guid", Filename: "filename"},
		//	},
		//},
	}
	// TODO добавить сохранение из json
	// TODO тесты для gin
	c.IndentedJSON(http.StatusOK, notes)
}

func GetActualNote(c *gin.Context) {
	var guidParam string = c.PostForm("guid")
	if len(guidParam) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Не указан параметр 'guid'"})
		return
	}

	var noteService, err = configureNoteService()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

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

	noteService, err := configureNoteService()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

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

	var noteService, err = configureNoteService()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

	err = noteService.DownNoteVersion(guidParam)
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

	var noteService, err = configureNoteService()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

	err = noteService.UpNoteVersion(guidParam)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"result": true})
}

func configureNoteService() (*entity.NoteServiceImpl, error) {
	var jdbcTemplate = db.JdbcTemplateImpl{
		DriverName: utils.MySqlDriverName,
		DbUrl:      utils.MySqlDataSource}

	endpoint := utils.MinioEndpoint
	accessKeyID := utils.MinioAccessKey
	secretAccessKey := utils.MinioSecretKey
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var minioService = filestorage.MinioServiceImpl{
		Client: minioClient}

	var noteService = entity.NoteServiceImpl{
		JdbcTemplate: &jdbcTemplate,
		MinioService: &minioService}

	return &noteService, nil
}
