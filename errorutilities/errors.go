package errorutilities

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidUUID     = NewInvalidUUIDError()
	ErrNotFound        = NewNotFoundError("resource")
	ErrFieldUnexported = NewError("field is unexported, can't get it's value")
	ErrDataNotStruct   = NewError("data isn't struct")
	ErrFieldNotInt     = NewError("field isn't type int, int8, int16, int32 or int64")
	ErrFieldNotUInt    = NewError("field isn't type uint, uint8, uint16, uint32 or uint64")
	ErrFieldNotString  = NewError("field isn't type string")
	ErrNotPointer      = NewError("destination should be pointer or else changes won't reflect on it")
)

type Error struct {
	Message string `json:"error"`
	Status  int    `json:"-"`
}

func (e Error) Error() string {
	return e.Message
}

func NewError(message string, status ...int) Error {
	finalStatus := http.StatusInternalServerError
	if len(status) > 0 {
		finalStatus = status[0]
	}

	return Error{
		Status:  finalStatus,
		Message: message,
	}
}

func NewErrorf(pattern string, args ...any) Error {
	return Error{
		Message: fmt.Sprintf(pattern, args...),
	}
}

func NewBadRequestError(message string) Error {
	return Error{
		Message: fmt.Sprintf("BadRequest: request doesn't fullfill the requirements: %s", message),
		Status:  http.StatusBadRequest,
	}
}

func NewNotFoundError(resource string) Error {
	return Error{
		Message: fmt.Sprintf("%s not found", resource),
		Status:  http.StatusNotFound,
	}
}

func NewInvalidUUIDError() Error {
	return Error{
		Message: "UUID needs to be formatted according to the UUID v4 convention.",
		Status:  http.StatusBadRequest,
	}
}

func NewValidationError(errors validator.ValidationErrors) Error {
	message := strings.Join(
		pie.Map(errors, func(e validator.FieldError) string {
			return e.Error()
		}),
		"\n",
	)

	return Error{
		Message: fmt.Sprintf("The following fields don't meet the validation requirements:\n%s", message),
		Status:  http.StatusBadRequest,
	}
}

func NewUnexportedFieldError(fieldName string) Error {
	return Error{
		Message: fmt.Sprintf("field is unexported, can't get it's value: %s", fieldName),
	}
}
