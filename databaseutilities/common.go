package databaseutilities

import "context"

type Database[T any] interface {
	Query(ctx context.Context, query string, args ...any) (T, error)
	Close()
}
