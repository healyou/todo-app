package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/pkg/errors"
)

type JdbcTemplateImpl struct {
	SqlDb *sql.DB
}

func (jdbcTemplate JdbcTemplateImpl) ExecuteInTransaction(
	txFunc func(context context.Context, DB *sql.Tx) error) error {

	ctx := context.Background()
	tx, err := jdbcTemplate.SqlDb.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "ошибка создания транзакции")
	}

	err = txFunc(ctx, tx)

	if err != nil {
		var txErr = tx.Rollback()
		if txErr != nil {
			log.Println(fmt.Printf("%+v", errors.Wrap(txErr, "ошибка отката транзакции")))
		}
		return errors.Wrap(err, "ошибка в транзакционной функции")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "ошибка в процессе фиксации транзакции")
	}

	return nil
}
