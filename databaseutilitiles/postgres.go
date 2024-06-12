package databaseutilitiles

import (
	"context"
	"errors"
	"reflect"

	"github.com/Originate/go-utilities/configutilities"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database[T any] interface {
	Query(ctx context.Context, query string, args ...any) (T, error)
	Close()
}

type PostgresConnection struct {
	db Database[pgx.Rows]
}

func NewPostgres(ctx context.Context, cfg configutilities.DatabasePostgresConfiguration) (*PostgresConnection, error) {
	db, err := pgxpool.New(ctx, buildPostgresDSN(cfg))
	if err != nil {
		return &PostgresConnection{}, err
	}

	return &PostgresConnection{
		db: db,
	}, nil
}

func (c *PostgresConnection) scan(rows pgx.Rows, dest interface{}) error {
	switch reflect.ValueOf(dest).Elem().Kind() {
	case reflect.Slice:
		return checkNoRowsInResultSetError(pgxscan.ScanAll(dest, rows))
	default:
		return checkNoRowsInResultSetError(pgxscan.ScanOne(dest, rows))
	}
}

func (c *PostgresConnection) Query(ctx context.Context, dest interface{}, query string, args ...any) error {
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return errors.New("destination should be pointer or else changes won't reflect on it")
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return checkNoRowsInResultSetError(err)
	}

	return c.scan(rows, dest)
}

func (c *PostgresConnection) Close() {
	c.db.Close()
}
