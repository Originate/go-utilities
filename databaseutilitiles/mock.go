package databaseutilitiles

import (
	"context"
	"reflect"

	"github.com/Originate/go-utilities/errorutilities"
	"github.com/pashagolub/pgxmock/v4"
)

type PostgresMockConn struct {
	Mock pgxmock.PgxPoolIface
}

func NewPostgresMock() (*PostgresMockConn, error) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		return &PostgresMockConn{}, err
	}

	return &PostgresMockConn{
		Mock: mock,
	}, nil
}

func (m *PostgresMockConn) Query(ctx context.Context, dest interface{}, query string, args ...any) error {
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return errorutilities.ErrNotPointer
	}

	rows, err := m.Mock.Query(ctx, query, args...)
	if err != nil {
		return checkNoRowsInResultSetError(err)
	}

	return scan(rows, dest)
}

func (m *PostgresMockConn) Close() {
	m.Mock.Close()
}
