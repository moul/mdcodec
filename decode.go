package mdcodec

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

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
	if len(lines) == 0 {
		return nil
	}

	for len(lines) > 0 {
		line := lines[0]
		lines = lines[1:]

		if strings.HasPrefix(line, "- **") {
			parts := strings.SplitN(line, "**:", 2)
			if len(parts) < 2 {
				continue
			}

			fieldName := strings.Trim(parts[0][4:], "* ")
			fieldValStr := strings.TrimSpace(parts[1])

			field := val.FieldByName(fieldName)
			if !field.IsValid() {
				return fmt.Errorf("no such field: %s in %v", fieldName, val.Type())
			}
			if !field.CanSet() {
				return fmt.Errorf("cannot set field: %s", fieldName)
			}

			switch field.Kind() {
			case reflect.Struct:
				if fieldValStr != "" {
					return errors.New("unexpected value for nested struct")
				}
				// Recursively parse the struct
				endIdx := findEndOfSubStruct(lines)
				if err := parseStruct(lines[:endIdx], field); err != nil {
					return err
				}
				lines = lines[endIdx:]
			default:
				if err := parseValue(fieldValStr, field); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func parseValue(s string, v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		v.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.SetBool(b)
	default:
		return fmt.Errorf("unsupported type: %s", v.Type().Name())
	}
	return nil
}

func findEndOfSubStruct(lines []string) int {
	indentCount := countLeadingSpaces(lines[0])

	for i, line := range lines {
		// We consider a line that has the same or fewer indents to be the end of the sub-struct.
		if countLeadingSpaces(line) <= indentCount {
			return i
		}
	}

	// If we don't find an earlier end, then the sub-structure extends to the end of the lines slice.
	return len(lines)
}

func countLeadingSpaces(s string) int {
	return len(s) - len(strings.TrimLeft(s, " "))
}
