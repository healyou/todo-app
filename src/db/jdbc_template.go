package db

import (
	"context"
	"database/sql"
)

// Выполнение операций с базой данных
type JdbcTemplate interface {
	ExecuteInTransaction(
		txFunc func(context context.Context, DB *sql.Tx) error) error
}
