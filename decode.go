package mdcodec

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Unmarshal will convert markdown content to a Go struct.
func Unmarshal(md string, v interface{}) error {
	lines := strings.Split(md, "\n")
	ptrValue := reflect.ValueOf(v)

	if ptrValue.Kind() != reflect.Ptr || ptrValue.IsNil() {
		return errors.New("expected pointer to a value")
	}

	val := ptrValue.Elem()
	return parseStruct(lines, val)
}

func parseStruct(lines []string, val reflect.Value) error {
	var i int
	for i = 0; i < len(lines); i++ {
		line := lines[i]

		if strings.HasPrefix(line, "# ") {
			// This denotes a new struct (or a title); we should move on
			continue
		}

		if strings.HasPrefix(line, "- **") {
			// This is a field
			parts := strings.SplitN(line, ":", 2)
			right := strings.TrimPrefix(parts[1], " ")
			fieldName := strings.TrimSuffix(strings.TrimPrefix(parts[0], "- **"), "**")

			// Check if we're dealing with a sub-struct or slice (indicated by a newline without a value after the colon)
			if len(parts) == 1 || right == "" {
				field, ok := fieldByName(val, fieldName)
				if !ok {
					return fmt.Errorf("no field with tag: %s", fieldName)
				}
				subLines, remainingLines := findEndOfSubStruct(lines[i+1:])
				if field.Kind() == reflect.Slice {
					// If it's a slice, we need to handle each element
					elemType := field.Type().Elem()
					slice := reflect.MakeSlice(field.Type(), 0, len(subLines))
					for _, itemStr := range subLines {
						item := reflect.New(elemType).Elem()
						// Assume each item in the slice is a single line for simplicity; adjust as needed
						if err := setFieldValue(item, fieldName, itemStr); err != nil {
							return err
						}
						slice = reflect.Append(slice, item)
					}
					field.Set(slice)
				} else if field.Kind() == reflect.Struct {
					// If it's a struct, recursively parse
					if err := parseStruct(subLines, field); err != nil {
						return err
					}
				}
				// Adjust the outer loop's index to skip over the lines we've just processed
				i += len(subLines)
				lines = remainingLines
				continue
			}

			fieldValueStr := right
			if err := setFieldValue(val, fieldName, fieldValueStr); err != nil {
				return err
			}
		}
	}
	return nil
}

func fieldByName(v reflect.Value, name string) (reflect.Value, bool) {
	field := v.FieldByName(name)
	if field.IsValid() {
		return field, true
	}
	return reflect.Value{}, false
}

func fieldByTag(v reflect.Value, tagValue string) (reflect.Value, bool) {
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("md")
		if tag == tagValue {
			return v.Field(i), true
		}
	}
	return reflect.Value{}, false
}

func findEndOfSubStruct(lines []string) (currentStruct, remaining []string) {
	indentCount := 0
	for i, line := range lines {
		if strings.HasPrefix(line, "- **") {
			if indentCount == 0 {
				return lines[:i], lines[i:]
			}
			indentCount--
		} else if strings.HasPrefix(line, "# ") {
			indentCount++
		}
	}
	return lines, nil
}

func setFieldValue(val reflect.Value, fieldName, fieldValueStr string) error {
	fieldType, exists := val.Type().FieldByName(fieldName)
	if !exists {
		return fmt.Errorf("field %s does not exist", fieldName)
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() || !field.CanSet() {
		return fmt.Errorf("field %s cannot be set", fieldName)
	}

	switch fieldType.Type.Kind() {
	case reflect.String:
		field.SetString(fieldValueStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(fieldValueStr, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	default:
		return fmt.Errorf("unsupported type %s", fieldType.Type.Kind())
	}

	return nil
}
