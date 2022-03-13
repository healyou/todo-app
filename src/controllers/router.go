package note_controller

import (
	"log"
	"todo/src/db"
	"todo/src/entity"
	"todo/src/filestorage"
	"todo/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func SetupRouter(setup func (*gin.Engine)) *gin.Engine {
	router := gin.Default()
	setup(router)
	router.POST("/notes/getActualNote", GetActualNote)
	router.POST("/notes/saveNote", SaveNote)
	router.POST("/notes/getUserNotes", GetUserNotes)
	router.POST("/notes/downNoteVersion", DownNoteVersion)
	router.POST("/notes/upNoteVersion", UpNoteVersion)
	return router
}

func SetupMiddleware(router *gin.Engine) {
	var noteService, err = configureNoteService()
	if err != nil {
		log.Fatalln(err)
		return
	}
	router.Use(ApiMiddleware(noteService))
}

func ApiMiddleware(noteService *entity.NoteServiceImpl) gin.HandlerFunc {
    // TODO эти штуки надо сделать как одиночки
	return func(c *gin.Context) {
        c.Set("noteService", noteService)
        c.Next()
    }
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