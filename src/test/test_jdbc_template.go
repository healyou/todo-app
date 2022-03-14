package test

import (
	"context"
	"database/sql"
)

type JdbcTemplateImplTest struct {
	DB      *sql.Tx
	context context.Context
}

func (jdbcTemplate JdbcTemplateImplTest) ExecuteInTransaction(txFunc func(context context.Context, DB *sql.Tx) error) error {
	err := txFunc(jdbcTemplate.context, jdbcTemplate.DB)

	if err != nil {
		return err
	}

	return nil
}

func (jdbcTemplate JdbcTemplateImplTest) InTransactionForSqlResult(
	txFunc func(context context.Context, DB *sql.Tx) (*sql.Result, error)) (*sql.Result, error) {

	result, err := txFunc(jdbcTemplate.context, jdbcTemplate.DB)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (jdbcTemplate JdbcTemplateImplTest) InTransactionForSqlRows(
	txFunc func(context context.Context, DB *sql.Tx) (*sql.Rows, error)) (*sql.Rows, error) {

	result, err := txFunc(jdbcTemplate.context, jdbcTemplate.DB)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (jdbcTemplate JdbcTemplateImplTest) InTransactionForSqlRow(
	txFunc func(context context.Context, DB *sql.Tx) (*sql.Row, error)) (*sql.Row, error) {

	result, err := txFunc(jdbcTemplate.context, jdbcTemplate.DB)

	if err != nil {
		return nil, err
	}

	return result, nil
}
