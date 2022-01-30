package db

import (
	"context"
	"database/sql"
	"log"
)

type JdbcTemplateImpl struct {
	DriverName string
	DbUrl      string
}

func (jdbcTemplate *JdbcTemplateImpl) InTransactionForSqlResult(
	txFunc func(context context.Context, DB *sql.Tx) (*sql.Result, error)) (*sql.Result, error) {

	db, err := sql.Open(jdbcTemplate.DriverName, jdbcTemplate.DbUrl)
	if err != nil {
		return nil, err
	}
	defer func(db *sql.DB) {
		err = db.Close()
	}(db)

	// Create a new context, and begin a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	result, err := txFunc(ctx, tx)

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return result, nil
}

func (jdbcTemplate *JdbcTemplateImpl) InTransactionForSqlRows(
	txFunc func(context context.Context, DB *sql.Tx) (*sql.Rows, error)) (*sql.Rows, error) {

	db, err := sql.Open(jdbcTemplate.DriverName, jdbcTemplate.DbUrl)
	if err != nil {
		return nil, err
	}
	defer func(db *sql.DB) {
		err = db.Close()
	}(db)

	// Create a new context, and begin a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	result, err := txFunc(ctx, tx)

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return result, nil
}

func (jdbcTemplate *JdbcTemplateImpl) InTransactionForSqlRow(
	txFunc func(context context.Context, DB *sql.Tx) (*sql.Row, error)) (*sql.Row, error) {

	db, err := sql.Open(jdbcTemplate.DriverName, jdbcTemplate.DbUrl)
	if err != nil {
		return nil, err
	}
	defer func(db *sql.DB) {
		err = db.Close()
	}(db)

	// Create a new context, and begin a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	result, err := txFunc(ctx, tx)

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return result, nil
}
