package mdcodec

import (
	"fmt"
	"reflect"
	"strings"
)

func Marshal(v interface{}) (string, error) {
	return marshalValue(reflect.ValueOf(v), ""), nil
}

func marshalValue(v reflect.Value, indent string) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", v.Float())
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool())
	case reflect.Slice, reflect.Array:
		var b strings.Builder
		for i := 0; i < v.Len(); i++ {
			b.WriteString(marshalValue(v.Index(i), indent))
		}
		return b.String()
	case reflect.Map:
		var b strings.Builder
		for _, key := range v.MapKeys() {
			b.WriteString(fmt.Sprintf("%s- **%s**: %s\n", indent, key.String(), marshalValue(v.MapIndex(key), "  "+indent)))
		}
		return b.String()
	case reflect.Struct:
		var b strings.Builder
		t := v.Type()

		isTopLevel := indent == ""
		if isTopLevel {
			titleField, hasTitle := findFieldByTag(v, "title")
			if hasTitle {
				b.WriteString(fmt.Sprintf("# %s (%s)\n\n", titleField.String(), t.Name()))
			} else {
				b.WriteString(fmt.Sprintf("# %s\n\n", t.Name()))
			}
		} else {
			b.WriteString("\n")
		}

		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			tags := parseFieldTags(field)
			if tags["title"] == "true" {
				continue // skip title field
			}
			nestedVal := marshalValue(v.Field(i), "  "+indent)
			// ugly hack to not append a space if the val is a complex object
			if strings.HasPrefix(nestedVal, "\n") {
				b.WriteString(fmt.Sprintf("%s- **%s**:%s\n", indent, field.Name, nestedVal))
			} else {
				b.WriteString(fmt.Sprintf("%s- **%s**: %s\n", indent, field.Name, nestedVal))
			}
		}
		return b.String()
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}
