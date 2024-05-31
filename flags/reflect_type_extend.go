package flags

import (
	"reflect"
)

var extends = extendMap{}

type ExtendType interface {
	// Get 用于根据传入的reflect.Value获取对应的字符串。
	//
	// 参数:
	//  - v - 代表要获取数据的reflect.Value。
	Get(v reflect.Value) (s string)

	// Set 用于通过反射设置值。
	//
	// 参数:
	//   - v: 要设置值的反射值。
	//   - s: 要设置的字符串值。
	Set(v reflect.Value, s string) (err error)

	// Type 值的类型字符串。
	Type() string
}

// Extend 为指定类型 T 注册自定义的解析和格式化函数。
//
//	parse 函数用于将字符串解析为 T 类型的值，format 函数用于将 T 类型的值格式化为字符串。
//
// 参数:
//   - parse: 一个函数，其输入为字符串，输出为 T 类型的值和可能的错误。用于将配置文件中的字符串值解析为实际的类型 T。
//   - format: 一个函数，其输入为 T 类型的值，输出为字符串。用于将类型 T 的值格式化为字符串，以便写入配置文件中。
func Extend[T any](parse func(string) (T, error), format func(T) string) {
	// 定义一个 setFunc，用于设置值到反射的 Value 中
	setFunc := func(v reflect.Value, s string) (err error) {
		// 如果 parse 函数不为空，则尝试使用 parse 函数解析字符串
		if parse != nil {
			// 尝试解析字符串，成功则设置到反射的 Value 中
			if r, e := parse(s); e == nil {
				v.Set(reflect.ValueOf(r))
			} else {
				// 解析失败，返回错误
				err = e
			}
		}
		return
	}

	// 定义一个 getFunc，用于从反射的 Value 中获取值并格式化为字符串
	getFunc := func(v reflect.Value) (s string) {
		// 如果 format 函数不为空，则尝试使用 format 函数格式化值为字符串
		if format != nil {
			r, _ := v.Interface().(T)
			s = format(r)
		}
		return
	}

	// 创建一个简单的类型，包含类型 T 的反射类型、设置函数和获取函数
	var x T
	st := &simpleType{typ: reflect.TypeOf(x), setFunc: setFunc, getFunc: getFunc}

	// 将该简单类型注册到 extends 映射中，用于后续的类型-函数绑定
	if extends == nil {
		// 如果 extends 映射为空，则初始化一个新的映射并插入当前类型-函数绑定
		extends = extendMap{st.typ: st}
	} else {
		// 如果 extends 映射已存在，则直接插入当前类型-函数绑定
		extends[st.typ] = st
	}
}

// GetExtend 通过反射类型获取对应的扩展类型
//
// 参数:
//   - t: 要获取扩展类型的反射类型
//
// 返回值:
//   - ExtendType: 与给定反射类型对应的扩展类型
func GetExtend(t reflect.Type) ExtendType { return extends[t] }

// IsExtend 函数用于判断指定类型是否为扩展类型。
//
// 参数:
//   - t reflect.Type - 需要判断类型的 reflect.Type 实例。
//
// 返回值:
//   - bool - 如果类型是扩展类型，则返回 true；否则返回 false。
func IsExtend(t reflect.Type) (yes bool) { _, yes = extends[t]; return }

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
