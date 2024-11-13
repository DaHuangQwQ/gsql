package reflect

import (
	"errors"
	"reflect"
)

// IterateFields
// It return error if entity's Kind is not [Struct] or [Pointer] to struct.
func IterateFields(entity any) (map[string]any, error) {
	// IterateFields(nil)
	if entity == nil {
		return nil, errors.New("entity is not a struct or a pointer to struct")
	}

	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)

	// IterateFields((*Struct)(nil))
	if val.IsZero() {
		return nil, errors.New("entity is not a struct or a pointer to struct")
	}

	// IterateFields(*Struct)
	for typ.Kind() == reflect.Pointer {
		// get val in pointer
		typ = typ.Elem()
		val = val.Elem()
	}

	// IterateFields(*int)
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("entity is not a struct or a pointer to struct")
	}

	numField := typ.NumField()

	fields := make(map[string]any, numField)

	for i := 0; i < numField; i++ {
		fieldType := typ.Field(i)
		fieldVal := val.Field(i)

		if fieldType.IsExported() {
			fields[fieldType.Name] = fieldVal.Interface()
		} else {
			fields[fieldType.Name] = reflect.Zero(fieldType.Type).Interface()
		}
	}
	return fields, nil
}

func SetField(entity any, field string, newValue any) error {
	val := reflect.ValueOf(entity)
	for val.Type().Kind() == reflect.Pointer {
		val = val.Elem()
	}

	if val.IsZero() {

	}

	fieldVal := val.FieldByName(field)
	if !fieldVal.CanSet() {
		return errors.New("cannot set field ")
	}
	fieldVal.Set(reflect.ValueOf(newValue))

	return nil
}
