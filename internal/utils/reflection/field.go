// Package reflection contains functions for working with reflection.

package reflection

import (
	"reflect"
)

// GetFieldName returns the name of the field pointer
// in the struct pointer is pointing to.
func GetFieldName(structPtr interface{}, fieldPtr interface{}) (string, bool) {
	if structPtr == nil {
		return "", false
	}

	if fieldPtr == nil {
		return "", false
	}

	structVal := reflect.ValueOf(structPtr)

	if structVal.Kind() != reflect.Pointer ||
		structVal.Elem().Kind() != reflect.Struct {
		return "", false
	}

	fieldVal := reflect.ValueOf(fieldPtr)

	if fieldVal.Kind() != reflect.Pointer {
		return "", false
	}

	structElem := structVal.Elem()

	// range shouldn't be used cause iteration over result of function
	//nolint:intrange
	for i := 0; i < structElem.NumField(); i++ {
		if fieldVal.Pointer() == structElem.Field(i).Addr().Pointer() {
			return structElem.Type().Field(i).Name, true
		}
	}

	return "", false
}
