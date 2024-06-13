package configutilities

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"
)

func HideSensitive[T any](cfg T) T {
	mask(reflect.ValueOf(cfg).Elem())

	return cfg
}

func mask(in reflect.Value) {
	element := in
	if in.Kind() == reflect.Ptr {
		element = in.Elem()
	}

	for i := 0; i < element.NumField(); i++ {
		field := element.Field(i)

		switch field.Kind() {
		case reflect.Struct, reflect.Ptr:
			mask(field)
		case reflect.String:
			raw := field.String()

			if field.CanSet() {
				if len(raw) > 4 {
					field.SetString(fmt.Sprintf("%s%s", strings.Repeat("*", 6), raw[len(raw)-2:]))
				} else {
					field.SetString(strings.Repeat("*", len(raw)))
				}
			} else {
				slog.Debug("field can't be set", "fieldName", field.Type().Field(i).Name)
			}
		}
	}
}
