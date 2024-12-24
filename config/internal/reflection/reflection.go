/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package reflection implements the functions, types, and interfaces for the module.
package reflection

import (
	"fmt"
	"reflect"
)

// FieldValueByType recursively traverses the fields of an any type and returns the value of the field with the specified type.
func FieldValueByType[T any](obj any) (T, error) {
	v := reflect.ValueOf(obj)
	t := reflect.TypeFor[T]()
	var zero T
	// If the object is a pointer, get the value it points to.
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// If the object is not a struct, return an error.
	if v.Kind() != reflect.Struct {
		return zero, fmt.Errorf("expected a struct, got %T", obj)
	}

	// Iterate over all fields of the struct.
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := field.Type()

		// If the field type matches the desired type, return the field value.
		if fieldType == t {
			return field.Interface().(T), nil
		}

		// If the field is a struct, recursively search for the desired type within it.
		if field.Kind() == reflect.Struct {
			value, err := FieldValueByType[T](field.Interface())
			if err == nil {
				return value, nil
			}
		}
	}

	// If the desired type is not found, return an error.
	return zero, fmt.Errorf("field of type %s not found in %T", t, obj)
}

func FieldPointByType[T any](obj any) (*T, error) {
	return FieldValueByType[*T](obj)
}
