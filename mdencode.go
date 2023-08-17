package mdencode

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Marshal(v interface{}) string {
	val := reflect.ValueOf(v)

	// We will use a StringBuilder for building our markdown representation
	var builder strings.Builder

	// We start by examining the type of the top-level struct
	writeStruct(&builder, val, "")

	return builder.String()
}

func writeStruct(builder *strings.Builder, val reflect.Value, prefix string) {
	typ := val.Type()

	// Begin the section with the name of the struct
	builder.WriteString(prefix + "# " + typ.Name() + "\n\n")

	for i := 0; i < val.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)

		switch fieldValue.Kind() {
		case reflect.String:
			builder.WriteString(fmt.Sprintf("- **%s**: %s\n", fieldType.Name, fieldValue.String()))
		case reflect.Int:
			builder.WriteString(fmt.Sprintf("- **%s**: %d\n", fieldType.Name, fieldValue.Int()))
		case reflect.Struct:
			builder.WriteString("\n")
			writeStruct(builder, fieldValue, prefix+"#")
		default:
			// You can handle other types or ignore them
		}
	}
	builder.WriteString("\n")
}

func Unmarshal(md string, v interface{}) error {
	lines := strings.Split(md, "\n")
	ptrValue := reflect.ValueOf(v)

	if ptrValue.Kind() != reflect.Ptr || ptrValue.IsNil() {
		return errors.New("expected pointer to a value")
	}

	val := ptrValue.Elem()

	_, err := parseStruct(lines, val, 0)
	return err
}

func parseStruct(lines []string, val reflect.Value, start int) (int, error) {
	// typ := val.Type()

	for i := start; i < len(lines); i++ {
		line := lines[i]

		// Check if it's a struct field line like "- **Name**: John"
		if strings.HasPrefix(line, "- **") {
			parts := strings.SplitN(line, "**:", 2)
			if len(parts) < 2 {
				continue
			}

			fieldName := strings.Trim(parts[0], "- ** ")
			fieldValueStr := strings.TrimSpace(parts[1])

			// Find the field in the struct
			field := val.FieldByName(fieldName)
			if !field.IsValid() {
				continue
			}

			switch field.Kind() {
			case reflect.String:
				field.SetString(fieldValueStr)
			case reflect.Int:
				intVal, err := strconv.Atoi(fieldValueStr)
				if err != nil {
					return i, err
				}
				field.SetInt(int64(intVal))
			case reflect.Struct:
				var err error
				i, err = parseStruct(lines, field, i+1)
				if err != nil {
					return i, err
				}
			default:
				// Handle other types or ignore them
			}
		} else if strings.HasPrefix(line, "# ") {
			// New struct detected, so we return
			return i - 1, nil
		}
	}

	return len(lines) - 1, nil
}
