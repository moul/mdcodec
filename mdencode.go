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
	writeStruct(&builder, val, "", "")

	return builder.String()
}

func writeStruct(builder *strings.Builder, val reflect.Value, name string, prefix string) {
	typ := val.Type()

	if name == "" {
		name = typ.Name()
	}

	// Begin the section with the name of the struct
	builder.WriteString(prefix + "# " + name + "\n\n")

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
			writeStruct(builder, fieldValue, fieldType.Name, prefix+"#")
		default:
			// You can handle other types or ignore them
		}
	}

	if prefix != "" {
		builder.WriteString("\n")
	}
}

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
	type context struct {
		prefix string
		val    reflect.Value
	}
	stack := []context{{prefix: "#", val: val}}

	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			prefix := strings.SplitN(line, " ", 2)[0]
			fieldName := strings.TrimPrefix(line, prefix+" ")

			if prefix == stack[len(stack)-1].prefix+"#" {
				parent := stack[len(stack)-1].val
				stack = append(stack, context{prefix: prefix, val: parent.FieldByName(fieldName)})
			} else {
				for len(stack) > 1 && stack[len(stack)-1].prefix != prefix {
					stack = stack[:len(stack)-1]
				}
			}
			continue
		}

		parts := strings.SplitN(line, "**: ", 2)
		if len(parts) == 2 {
			fieldName := strings.TrimPrefix(parts[0], "- **")
			field := stack[len(stack)-1].val.FieldByName(fieldName)
			if !field.IsValid() || !field.CanSet() {
				return fmt.Errorf("invalid field %s", fieldName)
			}
			switch field.Kind() {
			case reflect.String:
				field.SetString(parts[1])
			case reflect.Int:
				n, err := strconv.ParseInt(parts[1], 10, 64)
				if err != nil {
					return fmt.Errorf("not a number: %q", parts[1])
				}
				field.SetInt(n)
			}
		}
	}

	return nil
}
