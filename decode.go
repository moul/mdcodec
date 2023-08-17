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
