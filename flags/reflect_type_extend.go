package flags

import (
	"reflect"
)

var extends = extendMap{}

type ExtendType interface {
	Get(v reflect.Value) (s string)
	Set(v reflect.Value, s string) (err error)
	Type() string
}

func Extend[T any](parse func(string) (T, error), format func(T) string) {
	setFunc := func(v reflect.Value, s string) (err error) {
		if parse != nil {
			if r, e := parse(s); e == nil {
				v.Set(reflect.ValueOf(r))
			} else {
				err = e
			}
		}
		return
	}

	getFunc := func(v reflect.Value) (s string) {
		if format != nil {
			r, _ := v.Interface().(T)
			s = format(r)
		}
		return
	}

	var x T
	st := &simpleType{typ: reflect.TypeOf(x), setFunc: setFunc, getFunc: getFunc}
	if extends == nil {
		extends = extendMap{st.typ: st}
	} else {
		extends[st.typ] = st
	}
}

func GetExtend(t reflect.Type) ExtendType { return extends[t] }
func IsExtend(t reflect.Type) (yes bool)  { _, yes = extends[t]; return }

type extendMap map[reflect.Type]ExtendType

type simpleType struct {
	typ     reflect.Type
	setFunc func(reflect.Value, string) error
	getFunc func(reflect.Value) string
}

func (te *simpleType) Get(v reflect.Value) (s string) {
	if te != nil && te.getFunc != nil {
		s = te.getFunc(v)
	}
	return
}

func (te *simpleType) Set(v reflect.Value, s string) (err error) {
	if te != nil && te.setFunc != nil {
		if err = te.setFunc(v, s); err != nil {
			return
		}
	}
	return
}

func (te *simpleType) Type() string { return rType(te.typ, true) }
