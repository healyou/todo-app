package di

import (
	"log"
	"sync"
	"todo/src/db"
	"todo/src/entity"
	"todo/src/filestorage"
	"todo/src/utils"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var onceDi sync.Once
var instance *DependencyInjectionImpl = nil

/* Получить объект одиночку (перед использованием надо проинициализировать) */
func GetInstance() *DependencyInjectionImpl {	
	onceDi.Do(func() {
	   instance = new(DependencyInjectionImpl)
	})
	return instance
}

var onceInitDep sync.Once
/* Инициализация зависимостей приложения */
func InitDependency() {
	onceInitDep.Do(func() {
		var di = GetInstance()
		di.jdbcTemplate = db.JdbcTemplateImpl{
			DriverName: utils.MySqlDriverName,
			DbUrl:      utils.MySqlDataSource}

		endpoint := utils.MinioEndpoint
		accessKeyID := utils.MinioAccessKey
		secretAccessKey := utils.MinioSecretKey
		var minioClient, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: false,
		})
		if err != nil {
			log.Fatalln(err)
			return
		}

		di.minioClient = *minioClient

		di.minioService = filestorage.MinioServiceImpl{
			Client: minioClient}
		di.noteService = entity.NoteServiceImpl{
			JdbcTemplate: di.jdbcTemplate,
			MinioService: di.minioService}
	})
}

// TODO это надо отправить в тесты
/* Инициализация зависимостей приложения для тестов */
func InitForTest(jdbc db.JdbcTemplate, 
	noteService entity.NoteService,
	minioService filestorage.MinioService) {

	var di = GetInstance()
	di.minioService = minioService
	di.noteService = noteService
}

type DependencyInjectionImpl struct {
	// TODO правильно ли так наследование делать?
	DependencyInjection

	jdbcTemplate db.JdbcTemplate
	minioClient minio.Client
	// TODO что значит вернуть указатель, а что значит вернуть значение
	noteService  entity.NoteService
	minioService filestorage.MinioService
}

func (di DependencyInjectionImpl) GetNoteService() entity.NoteService {
	return di.noteService
}

func (di DependencyInjectionImpl) GetMinioService() filestorage.MinioService {
	return di.minioService
}
