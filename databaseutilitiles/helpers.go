package databaseutilitiles

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Originate/go-utilities/configutilities"
	"github.com/Originate/go-utilities/errorutilities"
	"github.com/jackc/pgx/v5"
)

func checkNoRowsInResultSetError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
		return errorutilities.ErrNotFound
	}

	return err
}

func buildPostgresDSN(config configutilities.DatabasePostgresConfiguration) string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		config.User,
		config.Pass,
		config.Host,
		config.Port,
		config.DatabaseName,
		config.SSLMode,
	)
}
