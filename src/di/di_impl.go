package di

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"todo/src/db"
	"todo/src/entity"
	"todo/src/environment"
	"todo/src/filestorage"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

var instance *DependencyInjection = nil
var hasInit = false

/* Получить объект одиночку (перед использованием надо проинициализировать) */
func GetInstance() DependencyInjection {
	if !hasInit {
		log.Println("инициализация одиночек")
		instance = initDependencyNotTest()
		hasInit = true
	}
	var depInj DependencyInjection = *instance
	return depInj
}

func SetDiFromTest(di *DependencyInjection) {
	if flag.Lookup("test.v") == nil {
		err := errors.New("устанавливать значение DependencyInjection руками можно только в тестах")
		log.Println(fmt.Printf("%+v", err))
		os.Exit(1)
	} else {
		instance = di
		hasInit = true
	}
}

/* Инициализация зависимостей приложения */
func initDependencyNotTest() *DependencyInjection {
	if flag.Lookup("test.v") != nil {
		err := errors.New("в тестах необходимо переопределить DependencyInjection")
		log.Println(fmt.Printf("%+v", err))
		os.Exit(1)
	}

	var di = new(DependencyInjectionImpl)

	/* Соединение с базой */
	var sqlDb, err = sql.Open(
		environment.GetEnvVariables().MySqlDriverName,
		environment.GetEnvVariables().MySqlDataSource)
	if err != nil {
		err = errors.Wrap(err, "Ошибка создания соединения с бд")
		log.Println(fmt.Printf("%+v", err))
		os.Exit(1)
	}
	err = sqlDb.Ping()
	if err != nil {
		err = errors.Wrap(err, "Ошибка проверки соединения с бд")
		log.Println(fmt.Printf("%+v", err))
		os.Exit(1)
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
		err = errors.Wrap(err, "ошибка подключения к minio")
		log.Println(fmt.Printf("%+v", err))
		os.Exit(1)
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
