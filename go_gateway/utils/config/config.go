package config

import (
	"fmt"
	"reflect"
)

func TypeNameOf(v interface{}) string {
	return TypeName(reflect.TypeOf(v))
}

func TypeName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		return fmt.Sprint("*", TypeName(t.Elem()))
	}
	return fmt.Sprintf("%s.%s", t.PkgPath(), t.String())
}
