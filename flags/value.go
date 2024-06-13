package flags

import (
	"reflect"
	"strings"
)

type Value struct {
	Ref  reflect.Value //引用对象
	typ  reflect.Type  //引用类型
	defs []string      //默认值字符串
}

func newValue(v reflect.Value) *Value { return &Value{Ref: v, typ: v.Type(), defs: rGet(v)} }

func (v *Value) Type() string { return rType(v.DirectType()) }

func (v *Value) Set(s string) (err error) { return v.SetString(o2s(s), false, false, true) }

func (v *Value) String() string {
	if len(v.defs) > 0 {
		if v.IsKind(reflect.Slice) {
			return "[" + strings.Join(v.defs, ",") + "]"
		} else {
			return v.defs[0]
		}
	}
	return ""
}

func (v *Value) SetString(args []string, reset, asDefault, refSync bool) (err error) {
	if refSync {
		for _, arg := range args {
			if err = rSet(v.Ref, arg, reset); err != nil {
				return
			}
		}
	}

	if asDefault {
		v.defs = args
	}
	return
}

func (v *Value) DirectType() reflect.Type      { return typeIndirect(v.typ) }
func (v *Value) IsKind(kind reflect.Kind) bool { return v.DirectType().Kind() == kind }
