package test

import (
	"context"
	"database/sql"
	"testing"
	"todo/src/di"
	"todo/src/entity"

	// "todo/src/entity"
	"todo/src/utils"
)

func ExecuteTestRollbackTransaction(
	t *testing.T, txFunc func()) {

	db, err := sql.Open(utils.MySqlDriverName, utils.MySqlDataSource)
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

	var testDiValue di.DependencyInjection = TestDependencyInjectionImpl{
		NoteServiceValue: noteService, MinioServiceValue: minioServiceImplTest}
	/* mocking global object function*/
	di.GetInstance = func() di.DependencyInjection {
		return testDiValue
	}

	txFunc()
}
