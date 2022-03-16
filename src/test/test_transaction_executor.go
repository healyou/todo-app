package test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"todo/src/di"
	"todo/src/entity"
	"todo/src/environment"

	// "todo/src/entity"
	"todo/src/utils"
)

func ExecuteTestRollbackTransaction(
	t *testing.T, txFunc func()) {

	/* грузим тестовые переменные */
	os.Setenv(utils.ProfileEnvName, "TEST")
	environment.GetEnvVariables()

	db, err := sql.Open(
		environment.GetEnvVariables().MySqlDriverName, 
		environment.GetEnvVariables().MySqlDataSource)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("an error '%s' was not create transaction", err)
	}
	defer func(db *sql.DB, tx *sql.Tx) {
		err = tx.Rollback()
		if err != nil {
			t.Errorf("rollback error: %s", err)
		}
		err = db.Close()
		if err != nil {
			t.Errorf("rollback error: %s", err)
		}
	}(db, tx)

	testJdbcTemplate := JdbcTemplateImplTest{DB: tx, context: ctx}

	minioServiceImplTest := MinioServiceImplTest{}
	var noteService = entity.NoteServiceImpl{
		JdbcTemplate: &testJdbcTemplate,
		MinioService: &minioServiceImplTest}

	/* mocking global object function*/
	depInj := new(di.DependencyInjectionImpl)
	depInj.Initialize(noteService, minioServiceImplTest)

	// TODO почитать про TestMain метод для инициализации тестов
	var value di.DependencyInjection = *depInj
	di.SetDiFromTest(&value)

	txFunc()
}
