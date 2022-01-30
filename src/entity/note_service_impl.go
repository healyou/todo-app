package entity

import (
	"context"
	"database/sql"
	"log"
	"todo/src/db"
)

type NoteServiceImpl struct {
	JdbcTemplate db.JdbcTemplate
}

func (service *NoteServiceImpl) Test() error {
	sqlFunc := func(context context.Context, DB *sql.Tx) (*sql.Result, error) {
		sqlCount := "select count(*) from note"
		var cnt int
		err := DB.QueryRowContext(context, sqlCount).Scan(&cnt)
		if err != nil {
			return nil, err
		}

		sqlInsert := "INSERT INTO note (note_guid, text, user_id)\n VALUES ('not guid1', 'note text1', 1)"
		result, err := DB.ExecContext(context, sqlInsert)
		if err != nil {
			return nil, err
		}

		sqlCount = "select count(*) from note"
		err = DB.QueryRowContext(context, sqlCount).Scan(&cnt)
		if err != nil {
			return nil, err
		}

		return &result, nil
	}

	result, err := service.JdbcTemplate.InTransactionForSqlResult(sqlFunc)
	if err != nil {
		return err
	}
	if result != nil {
		log.Println(result)
	}

	return nil
}
