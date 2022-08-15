package util

import (
	"fmt"
	"net/url"
	"reflect"
)

func PushToParameters(instance any, query *url.Values) {
	val := reflect.ValueOf(instance)
	for i := 0; i < val.Type().NumField(); i++ {
		str := fmt.Sprintf("%s", val.Field(i).Interface())
		if str != "" {
			query.Add(val.Type().Field(i).Tag.Get("json"), str)
		}
	}
}

func BoolToInt(in bool) int {
	var result int // default 0
	if in {
		result = 1
	}
	return result
}
