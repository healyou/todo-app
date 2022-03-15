package db

import (
	"context"
	"database/sql"
	"log"
)

type JdbcTemplateImpl struct {
	SqlDb *sql.DB
}

func (jdbcTemplate JdbcTemplateImpl) ExecuteInTransaction(
	txFunc func(context context.Context, DB *sql.Tx) error) error {

	ctx := context.Background()
	tx, err := jdbcTemplate.SqlDb.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	err = txFunc(ctx, tx)

	if err != nil {
		var txErr = tx.Rollback()
		if txErr != nil {
			log.Println(txErr)
		}
		log.Println(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (jdbcTemplate JdbcTemplateImpl) InTransactionForSqlResult(
	txFunc func(context context.Context, DB *sql.Tx) (*sql.Result, error)) (*sql.Result, error) {

	ctx := context.Background()
	tx, err := jdbcTemplate.SqlDb.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := txFunc(ctx, tx)

	if err != nil {
		var txErr = tx.Rollback()
		if txErr != nil {
			log.Println(txErr)
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func (jdbcTemplate JdbcTemplateImpl) InTransactionForSqlRows(
	txFunc func(context context.Context, DB *sql.Tx) (*sql.Rows, error)) (*sql.Rows, error) {

	ctx := context.Background()
	tx, err := jdbcTemplate.SqlDb.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := txFunc(ctx, tx)

	if err != nil {
		var txErr = tx.Rollback()
		if txErr != nil {
			log.Println(txErr)
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func (jdbcTemplate JdbcTemplateImpl) InTransactionForSqlRow(
	txFunc func(context context.Context, DB *sql.Tx) (*sql.Row, error)) (*sql.Row, error) {

	ctx := context.Background()
	tx, err := jdbcTemplate.SqlDb.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
	}

	result, err := txFunc(ctx, tx)

	if err != nil {
		var txErr = tx.Rollback()
		if txErr != nil {
			log.Println(txErr)
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}
