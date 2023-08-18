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
	return parseStruct(lines, val, "")
}

func parseStruct(lines []string, val reflect.Value, indent string) error {
	var i int
	for i = 0; i < len(lines); i++ {
		line := lines[i]

		// Handling titles and types
		if strings.HasPrefix(line, indent+"# ") {
			// Parsing title and potential type
			titleContent := strings.TrimPrefix(line, indent+"# ")
			parts := strings.Split(titleContent, " (")
			title := parts[0]
			var typeName string
			if len(parts) > 1 {
				typeName = strings.TrimSuffix(parts[1], ")")
			} else {
				typeName = title
			}

			// If the type name matches the current struct's type name, continue parsing it
			// Else, it's a sign to end the parsing of the current struct
			if typeName != val.Type().Name() {
				break
			}

			// Setting the title field value if it exists
			titleField, hasTitle := fieldByTag(val, "title")
			if hasTitle && titleField.Kind() == reflect.String {
				titleField.SetString(title)
			}
			continue
		}

		// Handling fields and potential sub-structures
		if strings.HasPrefix(line, indent+"- **") {
			parts := strings.SplitN(line, ":", 2)
			fieldName := strings.TrimSuffix(strings.TrimPrefix(parts[0], indent+"- **"), "**")

			if len(parts) == 1 || strings.TrimSpace(parts[1]) == "" {
				field, ok := fieldByTag(val, fieldName)
				if !ok {
					field, ok = fieldByName(val, fieldName)
				}
				if !ok {
					return fmt.Errorf("no field with tag or name: %s", fieldName)
				}

				// For nested structs or slices
				subLines, remainingLines := findEndOfSubStruct(lines[i+1:], indent+"  ")

				if field.Kind() == reflect.Slice {
					elemType := field.Type().Elem()
					slice := reflect.MakeSlice(field.Type(), 0, len(subLines))
					for _, itemStr := range subLines {
						item := reflect.New(elemType).Elem()
						if err := setFieldValue(item, fieldName, itemStr); err != nil {
							return err
						}
						slice = reflect.Append(slice, item)
					}
					field.Set(slice)
				} else if field.Kind() == reflect.Struct {
					if err := parseStruct(subLines, field, indent+"  "); err != nil {
						return err
					}
				}

				i += len(subLines)
				lines = remainingLines
				continue
			}

			fieldValueStr := strings.TrimSpace(parts[1])
			if err := setFieldValue(val, fieldName, fieldValueStr); err != nil {
				return err
			}
		}
	}
	return nil
}

func findEndOfSubStruct(lines []string, indent string) (currentStruct, remaining []string) {
	for i, line := range lines {
		if !strings.HasPrefix(line, indent) {
			return lines[:i], lines[i:]
		}
	}
	return lines, nil
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
