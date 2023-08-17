package mdcodec

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Marshal(v interface{}) (string, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return "", errors.New("expected pointer to a value")
	}

	val = val.Elem()
	return marshalStruct(val, ""), nil
}

func marshalStruct(val reflect.Value, indent string) string {
	var result strings.Builder

	// Extract name using the md:"title" tag
	nameField := findFieldNameByTag(val.Type(), "title")
	if nameField != "" {
		name := val.FieldByName(nameField).String()
		if indent == "" { // Only add type for the top-level structure
			result.WriteString(fmt.Sprintf("# %s (%s)\n", name, val.Type().Name()))
		}
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)
		if fieldType.Name == nameField {
			continue
		}

		switch field.Kind() {
		case reflect.Struct:
			result.WriteString(fmt.Sprintf("%s- **%s**:\n", indent, fieldType.Name))
			result.WriteString(marshalStruct(field, indent+"  "))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result.WriteString(fmt.Sprintf("%s- **%s**: %d\n", indent, fieldType.Name, field.Int()))
		case reflect.Float32, reflect.Float64:
			result.WriteString(fmt.Sprintf("%s- **%s**: %f\n", indent, fieldType.Name, field.Float()))
		case reflect.Bool:
			result.WriteString(fmt.Sprintf("%s- **%s**: %t\n", indent, fieldType.Name, field.Bool()))
		case reflect.String:
			result.WriteString(fmt.Sprintf("%s- **%s**: %s\n", indent, fieldType.Name, field.String()))
		default:
			result.WriteString(fmt.Sprintf("%s- **%s**: %v\n", indent, fieldType.Name, field.Interface()))
		}
	}

	return result.String()
}
