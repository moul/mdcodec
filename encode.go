package mdcodec

import (
	"fmt"
	"reflect"
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
