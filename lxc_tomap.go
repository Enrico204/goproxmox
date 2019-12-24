package goproxmox

import (
	"fmt"
	"reflect"
	"strings"
)

func (lxc *LXC) ToMap() map[string]string {
	postVars := map[string]string{}

	val := reflect.ValueOf(lxc).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		name := strings.ToLower(typeField.Name)
		if tag.Get("override") != "" {
			name = tag.Get("override")
		}
		if strings.Contains(tag.Get("json"), "omitempty") && valueField.IsZero() {
			continue
		}

		if valueField.Kind() == reflect.Slice {
			for i := 0; i < valueField.Len(); i++ {
				item := valueField.Index(i)
				if item.Kind() == reflect.TypeOf(VBaseNICSettings{}).Kind() {
					v := reflect.ValueOf(item.Interface()).MethodByName("ToProxmoxString").Call([]reflect.Value{reflect.ValueOf("lxc")})
					postVars[fmt.Sprintf("%s%d", name, i)] = v[0].String()
				} else if item.Kind() == reflect.String && !item.IsZero() {
					postVars[fmt.Sprintf("%s%d", name, i)] = item.String()
				}
			}
		} else if valueField.Kind() == reflect.Ptr && !valueField.IsNil() {
			elem := valueField.Elem()
			postVars[name] = fmt.Sprint(elem.Interface())
		} else {
			postVars[name] = fmt.Sprint(valueField.Interface())
		}
	}
	return postVars
}
