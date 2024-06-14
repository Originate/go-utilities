package databaseutilitiles

import (
	"context"
	"errors"
	"reflect"

	"github.com/Originate/go-utilities/configutilities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

func (c *PostgresConnection) Query(ctx context.Context, dest interface{}, query string, args ...any) error {
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return errors.New("destination should be pointer or else changes won't reflect on it")
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return checkNoRowsInResultSetError(err)
	}

	return scan(rows, dest)
}

func (c *PostgresConnection) Close() {
	c.db.Close()
}
