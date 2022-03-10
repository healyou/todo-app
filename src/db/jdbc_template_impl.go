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

func (jdbcTemplate *JdbcTemplateImpl) ExecuteInTransaction(
	txFunc func(context context.Context, DB *sql.Tx) error) error {

	db, err := sql.Open(jdbcTemplate.DriverName, jdbcTemplate.DbUrl)
	if err != nil {
		log.Println(err)
		return err
	}
	defer func(db *sql.DB) {
		err = db.Close()
	}(db)

	// Create a new context, and begin a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
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
