package goproxmox

import (
	"fmt"
	"reflect"
	"strings"
)

func (net *Network) ToMap() map[string]string {
	postVars := map[string]string{}

	val := reflect.ValueOf(net).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		if valueField.Kind() == reflect.String {
			postVars[strings.ToLower(typeField.Name)] = fmt.Sprint(valueField.Interface())
		} else if valueField.Kind() == reflect.Ptr && !valueField.IsNil() {
			elem := valueField.Elem()
			postVars[strings.ToLower(typeField.Name)] = fmt.Sprint(elem.Interface())
		}
	}
	return postVars
}
