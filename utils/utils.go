package utils

import "reflect"

func IsArray(data interface{}) bool {
	value := reflect.ValueOf(data)
	return value.Kind() == reflect.Slice || value.Kind() == reflect.Array
}
