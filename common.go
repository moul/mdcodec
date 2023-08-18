package mdcodec

import (
	"reflect"
	"strings"
)

func parseFieldTags(field reflect.StructField) map[string]string {
	tag := field.Tag.Get("md")
	parts := strings.Split(tag, ",")

	tags := make(map[string]string)
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		key := kv[0]
		value := "true" // default value for tags like "title"
		if len(kv) > 1 {
			value = kv[1]
		}
		tags[key] = value
	}
	return tags
}

func findFieldByTag(v reflect.Value, tagKey string) (reflect.Value, bool) {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tags := parseFieldTags(field)
		if _, exists := tags[tagKey]; exists {
			return v.Field(i), true
		}
	}
	return reflect.Value{}, false
}
