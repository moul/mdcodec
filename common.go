package mdcodec

import (
	"fmt"
	"reflect"
)

func setFieldValue(field reflect.Value, value string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		var intValue int
		fmt.Sscan(value, &intValue)
		field.SetInt(int64(intValue))
	}
}

func findFieldNameByTag(t reflect.Type, tagValue string) string {
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("md")
		if tag == tagValue {
			return t.Field(i).Name
		}
	}
	return ""
}
