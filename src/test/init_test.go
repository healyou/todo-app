package test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"todo/src/di"
	"todo/src/entity"
	"todo/src/environment"
	"todo/src/utils"
)

var db *sql.DB

/* 1 раз перед всеми тестами */
func TestMain(m *testing.M) {
	/* грузим тестовые переменные */
	loadTestEnv()
	/* подключаемся к базе */
	db = openMysql()

	/* Запуск тестов */
	exitVal := m.Run()

	/* отключаемся от бд */
	closeDb(db)

	os.Exit(exitVal)
}

func InitIntegrationTest(t *testing.T) func(t *testing.T) {
	/* Создаём транзакцию */
	ctx := context.Background()
	log.Println("создание транзакции для теста")
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("an error %s was not create transaction", err)
	}

	/* мокаем глобальные объекты */
	testJdbcTemplate := JdbcTemplateImplTest{DB: tx, context: ctx}
	minioServiceImplTest := MinioServiceImplTest{}
	var noteService = entity.NoteServiceImpl{
		JdbcTemplate: &testJdbcTemplate,
		MinioService: &minioServiceImplTest}

	depInj := new(di.DependencyInjectionImpl)
	depInj.Initialize(noteService, minioServiceImplTest)
	var value di.DependencyInjection = *depInj

	di.SetDiFromTest(&value)

	/* Функция завершения теста (очистка транзакции) */
	return func(t *testing.T) {
		log.Println("откат транзакции для теста")
		err = tx.Rollback()
		if err != nil {
			t.Errorf("rollback error: %s", err)
		}
	}
}

func loadTestEnv() {
	os.Setenv(utils.ProfileEnvName, "TEST")
	environment.GetEnvVariables()
}

func openMysql() *sql.DB {
	log.Println("открытие соединения с бд")
	db, err := sql.Open(
		environment.GetEnvVariables().MySqlDriverName, 
		environment.GetEnvVariables().MySqlDataSource)
	if err != nil {
		log.Fatalln("an error was not expected when opening a stub database connection", err)
	}
	return db
}

func closeDb(db *sql.DB) {
	log.Println("закрытие соединения с бд")
	err := db.Close()
	if err != nil {
		log.Fatalln("close db error:", err)
	}
}