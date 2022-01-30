package test

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"todo/src/entity"
)

func TestGgWp(t *testing.T) {
	txFunc := func(testJdbcTemplate JdbcTemplateImplTest) {
		var noteService = entity.NoteServiceImpl{JdbcTemplate: &testJdbcTemplate}
		err := noteService.Test()
		if err != nil {
			t.Errorf("error was not expected while test method: %s", err)
		}
	}
	ExecuteTestRollbackTransaction(t, txFunc)
}
