package db

import (
	"context"
	"database/sql"
)

// JdbcTemplate (SQL Go database connection) is a wrapper for SQL database handler ( can be *sql.DB or *sql.Tx)
// It should be able to work with all SQL data that follows SQL standard.
type JdbcTemplate interface {
	ExecuteInTransaction(
		txFunc func(context context.Context, DB *sql.Tx) error) error

	InTransactionForSqlResult(
		txFunc func(context context.Context, DB *sql.Tx) (*sql.Result, error)) (*sql.Result, error)
	InTransactionForSqlRows(
		txFunc func(context context.Context, DB *sql.Tx) (*sql.Rows, error)) (*sql.Rows, error)
	InTransactionForSqlRow(
		txFunc func(context context.Context, DB *sql.Tx) (*sql.Row, error)) (*sql.Row, error)
}
