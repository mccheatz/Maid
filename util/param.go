package util

import (
	"fmt"
	"net/url"
	"reflect"
)

func PushToParameters(instance any, query *url.Values) {
	val := reflect.ValueOf(instance)
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		query.Add(field.Tag.Get("json"), fmt.Sprintf("%s", val.Field(i).Interface()))
	}
}
