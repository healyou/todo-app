package di

import (
	"database/sql"
	"errors"
	"flag"
	"log"
	"todo/src/db"
	"todo/src/entity"
	"todo/src/environment"
	"todo/src/filestorage"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var instance *DependencyInjection = nil
var hasInit = false

/* Получить объект одиночку (перед использованием надо проинициализировать) */
func GetInstance() DependencyInjection {
	if !hasInit {
		instance = initDependencyNotTest()
		hasInit = true
	}
	var depInj DependencyInjection = *instance
	return depInj
}

func SetDiFromTest(di *DependencyInjection) {
	if flag.Lookup("test.v") == nil {
		log.Fatalln(errors.New("устанавливать значение DependencyInjection руками можно только в тестах"))
	} else {
		instance = di
		hasInit = true
	}
}

/* Инициализация зависимостей приложения */
func initDependencyNotTest() *DependencyInjection {
	if flag.Lookup("test.v") != nil {
		log.Fatalln(errors.New("в тестах необходимо переопределить DependencyInjection"))
	}

	var di = new(DependencyInjectionImpl)

	/* Соединение с базой */
	var sqlDb, err = sql.Open(
		environment.GetEnvVariables().MySqlDriverName,
		environment.GetEnvVariables().MySqlDataSource)
	if err != nil {
		log.Println("Ошибка создания соединения с бд")
		log.Fatalln(err)
	}
	err = sqlDb.Ping()
	if err != nil {
		log.Println("Ошибка проверки соединения с бд")
		log.Fatalln(err)
	}

	/* Minio client */
	endpoint := environment.GetEnvVariables().MinioEndpoint
	accessKeyID := environment.GetEnvVariables().MinioAccessKey
	secretAccessKey := environment.GetEnvVariables().MinioSecretKey
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Println("ошибка подключения к minio")
		log.Fatalln(err)
	}

	/* Устанавливаем значения */
	jdbcTemplate := db.JdbcTemplateImpl{SqlDb: sqlDb}
	minioService := filestorage.MinioServiceImpl{
		Client: minioClient}
	noteService := entity.NoteServiceImpl{
		JdbcTemplate: jdbcTemplate,
		MinioService: di.minioService}
	di.Initialize(noteService, minioService)

	var depInj DependencyInjection = *di
	return &depInj
}

type DependencyInjectionImpl struct {
	noteService  entity.NoteService
	minioService filestorage.MinioService
}

func (depInj *DependencyInjectionImpl) Initialize(
	noteService entity.NoteService,
	minioService filestorage.MinioService) {

	depInj.noteService = noteService
	depInj.minioService = minioService
}

func (depInj DependencyInjectionImpl) GetNoteService() entity.NoteService {
	return depInj.noteService
}

func (depInj DependencyInjectionImpl) GetMinioService() filestorage.MinioService {
	return depInj.minioService
}
