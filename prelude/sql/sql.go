package sql

import (
	"context"
	"database/sql"
)

func Query[T any](db *sql.DB, callback func(rows *sql.Rows, v *T) error, query string, args ...any) ([]T, error) {
	return QueryContext(context.Background(), db, callback, query, args...)
}

func QueryContext[T any](ctx context.Context, db *sql.DB, scanner func(rows *sql.Rows, v *T) error, query string, args ...any) ([]T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var vs []T
	for rows.Next() {
		var v T
		err = scanner(rows, &v) //should be reference?
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	return vs, rows.Err()
}

func QueryRow(db *sql.DB, scanner func(row *sql.Row) error, query string, args ...any) error {
	return QueryRowContext(context.Background(), db, scanner, query, args...)
}

func QueryRowContext(ctx context.Context, db *sql.DB, scanner func(row *sql.Row) error, query string, args ...any) error {
	return scanner(db.QueryRowContext(ctx, query, args...))
}

func Exec(db *sql.DB, query string, args ...any) error {
	return ExecContext(context.Background(), db, query, args...)
}

func ExecContext(ctx context.Context, db *sql.DB, query string, args ...any) error {
	_, err := db.ExecContext(ctx, query, args...)
	return err
}

/*
package psql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type C interface {
	Close(ctx context.Context) error
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Q
	P
}

type T interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	P
	Q
}

type P interface {
	Prepare(ctx context.Context, name string, sql string) (sd *pgconn.StatementDescription, err error)
}

type Q interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func Query[T any](q Q, query string, scanner func(r pgx.Rows, v *T) error, args ...any) ([]T, error) {
	return QueryContext(context.Background(), q, query, scanner, args...)
}

func QueryRow(q Q, query string, scanner func(r pgx.Row) error, args ...any) error {
	return QueryRowContext(context.Background(), q, query, scanner, args...)
}

func Exec(q Q, query string, args ...any) error {

	return ExecContext(context.Background(), q, query, args...)
}

func QueryContext[T any](ctx context.Context, q Q, query string, scanner func(r pgx.Rows, v *T) error, args ...any) ([]T, error) {
	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var vs []T
	for rows.Next() {
		var v T
		err = scanner(rows, &v)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	return vs, rows.Err()
}

func QueryRowContext(ctx context.Context, q Q, query string, scanner func(r pgx.Row) error, args ...any) error {
	return scanner(q.QueryRow(ctx, query, args...))
}

func ExecContext(ctx context.Context, q Q, query string, args ...any) error {
	_, err := q.Exec(ctx, query, args...)
	return err
}

*/
