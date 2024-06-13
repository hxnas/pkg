package flags

import (
	"fmt"
	"reflect"
	"strings"
)

type BindOption struct {
	Flag string
	Env  string
}

// StructBindE 通过反射将命令行参数绑定到结构体指针上。
//
// 参数:
//   - structPtr: 任意类型的结构体指针，将以此为模板解析命令行参数。
func StructBind(flag *FlagSet, structPtr any, prefix *Prefix) (err error) {
	var fields []*FlagField
	if fields, err = ParseStruct(structPtr, prefix); err != nil {
		return
	}

	for _, field := range fields {
		if err = field.applyDefault(); err != nil {
			return
		}

		usage := field.Usage
		if usage == "" {
			usage = field.Field.Name
		}

		if len(field.Env) > 0 {
			usage += fmt.Sprintf(" (env: %s)", strings.Join(field.Env, ", "))
		}

		// 创建并配置命令行参数项
		item := flag.VarPF(field.Value, field.Name, field.Shorthand, usage)
		item.Deprecated = field.Deprecated               // 设置字段的弃用信息
		item.ShorthandDeprecated = field.ShortDeprecated // 设置字段的简写弃用信息

		if fv := reflect.Indirect(field.Value.Ref); fv.Kind() == reflect.Bool {
			item.NoOptDefVal = "true"
		}
	}

	return nil
}

// FieldsWalk 打印结构体指针的字段信息
func FieldsWalk(structPtr any, prefix *Prefix, walFn func(field *FlagField, max int)) error {
	fields, err := ParseStruct(structPtr, prefix)
	if err != nil {
		return err
	}

	max := 0
	for _, f := range fields {
		if l := len(f.Field.Name); l > max {
			max = l
		}
	}

	// 格式化并打印每个字段的名称和值
	for _, f := range fields {
		walFn(f, max)
	}

	return nil
}

// StructToArgs 函数将结构体指针转换为命令行参数字符串切片。
//
// 参数:
//   - structPtr 任意类型的结构体指针，函数将通过反射解析其字段和值。
//
// 返回值 :
//
//   - args 是包含结构体字段及其值的字符串切片，格式为 "--字段名 值"。
func FieldsToArgs(structPtr any) (args []string, err error) {
	err = FieldsWalk(structPtr, nil, func(field *FlagField, _ int) {
		for _, s := range rGet(field.Value.Ref) {
			args = append(args, "--"+field.Name, s)
		}
	})
	return
}
