package lolp

import (
	"errors"
	"fmt"
	"reflect"
)

// SetField sets struct from key/value
func SetField(obj interface{}, field string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(field)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("no such field: %s in obj", field)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value", field)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("provided value type didn't match obj field type")
	}

	if reflect.TypeOf(val).String() != "" {
		structFieldValue.Set(val)
	}
	return nil
}
