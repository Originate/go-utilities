package databaseutilities

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Originate/go-utilities/configutilities"
	"github.com/Originate/go-utilities/errorutilities"
	"github.com/Originate/go-utilities/utils"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var invalidStringRegex = regexp.MustCompile("<([a-zA-Z0-9]) value>")
var updateQueryTemplate = "UPDATE %s SET %s WHERE id = $%d RETURNING *"

func scan(rows pgx.Rows, dest interface{}) error {
	switch reflect.ValueOf(dest).Elem().Kind() {
	case reflect.Slice:
		return checkNoRowsInResultSetError(pgxscan.ScanAll(dest, rows))
	default:
		return checkNoRowsInResultSetError(pgxscan.ScanOne(dest, rows))
	}
}

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

type UpdateQueryResult struct {
	Query string
	Args  []string
}

func BuildUpdateWithValues(tableName string, id string, data any) (UpdateQueryResult, error) {
	values := reflect.ValueOf(data)

	if values.Kind() != reflect.Struct {
		return UpdateQueryResult{}, errorutilities.ErrDataNotStruct
	}

	var args []string
	var statements []string
	for i := 0; i < values.NumField(); i++ {
		fieldValue := values.Field(i)

		if !fieldValue.IsZero() {
			fieldName := fieldValue.Elem().Type().Field(i).Name
			dbFieldName := utils.GetTagValue(values, i, "db")

			if dbFieldName == "" {
				dbFieldName = utils.ToSnakeCase(fieldName)
			}

			statements = append(
				statements,
				fmt.Sprintf("%s = $%d", dbFieldName, len(args)+1),
			)

			switch fieldValue.Type().Name() {
			case "string":
				stringValue := fieldValue.String()

				if invalidStringRegex.Match([]byte(stringValue)) {
					return UpdateQueryResult{}, errorutilities.ErrFieldNotString
				}

				args = append(args, stringValue)
			case "int", "int8", "int16", "int32", "int64":
				if !fieldValue.CanInt() {
					return UpdateQueryResult{}, errorutilities.ErrFieldNotInt
				}

				args = append(args, strconv.Itoa(int(fieldValue.Int())))
			case "uint", "uint8", "uint16", "uint32", "uint64":
				if !fieldValue.CanUint() {
					return UpdateQueryResult{}, errorutilities.ErrFieldNotUInt
				}

				args = append(args, strconv.FormatUint(fieldValue.Uint(), 10))
			case "UUID":
				if !fieldValue.CanInterface() {
					return UpdateQueryResult{}, errorutilities.NewUnexportedFieldError(fieldName)
				}

				value, ok := fieldValue.Interface().(uuid.UUID)

				if ok {
					args = append(args, value.String())
				} else {
					return UpdateQueryResult{}, errorutilities.NewError(
						fmt.Sprintf("can't cast value to uuid.UUID type: %v", reflect.TypeOf(value).Name()),
					)
				}
			case "bool":
				args = append(args, fmt.Sprintf("%v", fieldValue.Bool()))
			default:
				log.Printf("field type not supported: %v", fieldValue.Type().Name())
			}
		}
	}

	return UpdateQueryResult{
		Query: fmt.Sprintf(
			updateQueryTemplate,
			tableName,
			strings.Join(statements, ", "),
			len(args)+1,
		),
		Args: append(args, id),
	}, nil
}
