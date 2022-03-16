package di

import (
	"database/sql"
	"errors"
	"flag"
	"log"
	"sync"
	"todo/src/db"
	"todo/src/entity"
	"todo/src/environment"
	"todo/src/filestorage"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

/* public reference to function (mock in tests) */
var GetInstance func() DependencyInjection = getInstanceFunction

var onceDi sync.Once
var instance *dependencyInjectionImpl = nil

/* Получить объект одиночку (перед использованием надо проинициализировать) */
func getInstanceFunction() DependencyInjection {
	onceDi.Do(func() {
		checkNoteTestModeForCreateDi()
		instance = new(dependencyInjectionImpl)
		initDependency(instance)
	})
	var depInj DependencyInjection = *instance
	return depInj
}

func checkNoteTestModeForCreateDi() {
	if flag.Lookup("test.v") != nil {
		log.Fatalln(errors.New("в тестах необходимо переопределять создание глобальных объектов"))
	} else {
		log.Println("создание глобальных объектов")
	}
}

/* Инициализация зависимостей приложения */
func initDependency(di *dependencyInjectionImpl) {
	/* Соединение с базой */
	var sqlDb, err = sql.Open(
		environment.GetEnvVariables().MySqlDriverName, 
		environment.GetEnvVariables().MySqlDataSource)
	if err != nil {
		log.Println("Ошибка создания соединения с бд")
		log.Fatalln(err)
		return
	}
	err = sqlDb.Ping()
	if err != nil {
		log.Println("Ошибка проверки соединения с бд")
		log.Fatalln(err)
		return
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
		return
	}

	/* Устанавливаем значения */
	di.sqlDb = sqlDb
	di.jdbcTemplate = db.JdbcTemplateImpl{SqlDb: sqlDb}
	di.minioClient = *minioClient

	di.minioService = filestorage.MinioServiceImpl{
		Client: minioClient}
	di.noteService = entity.NoteServiceImpl{
		JdbcTemplate: di.jdbcTemplate,
		MinioService: di.minioService}
}

type dependencyInjectionImpl struct {
	sqlDb        *sql.DB
	jdbcTemplate db.JdbcTemplate
	minioClient  minio.Client
	noteService  entity.NoteService
	minioService filestorage.MinioService
}

func (di dependencyInjectionImpl) GetNoteService() entity.NoteService {
	return di.noteService
}

func (di dependencyInjectionImpl) GetMinioService() filestorage.MinioService {
	return di.minioService
}
