package lod

import (
	"fmt"
	"reflect"
	"strconv"
)

// Itoa converts a number or bool to string
func Itoa[T Number | ~bool](in T) string {
	var z T
	if in == z {
		return ""
	}
	switch v := reflect.ValueOf(in); v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Bool:
		return "true"
	default:
		return fmt.Sprintf("%v", in)
	}
}

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Uint interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Float interface {
	~float32 | ~float64
}

type Number interface {
	Int | Uint | Float
}
