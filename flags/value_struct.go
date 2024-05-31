package flags

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// ParseStruct 是一个解析结构体的函数。
//   - 它接受一个指向结构体的指针作为参数，尝试使用默认配置解析传入的结构体。
//   - 如果解析失败，将错误信息打印到标准错误输出，并退出程序。
//
// 参数:
//   - structPtr any - 指向要解析的结构体的指针。
func ParseStruct(structPtr any) {
	if err := ParseStructE(Default(), structPtr, os.Args[1:]...); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// ParseStructE 是一个解析结构体参数的函数。
//   - 它首先尝试将参数绑定到指定的结构体指针上，然后解析额外的命令行参数。
//
// 参数:
//   - flag - 指向 FlagSet 的指针，用于处理命令行参数。
//   - structPtr - 任意类型的指针，预期为一个结构体指针，将从中解析参数。
//   - args - 可变数量的字符串参数，表示额外的命令行参数需要被解析。
func ParseStructE(flag *FlagSet, structPtr any, args ...string) (err error) {
	if err = StructBindE(structPtr); err != nil {
		return
	}

	if err = ParseFlag(flagSet(flag), args); err != nil {
		return
	}

	return
}

// StructBind 将命令行参数绑定到指定的结构体。该函数首先尝试进行绑定操作，如果操作失败，则会打印错误信息并退出程序。
//
// 参数:
//   - structPtr - 指向要绑定数据的结构体的指针。
//   - flags - 可选的FlagSet数组，用于定义额外的标志处理规则。
func StructBind(structPtr any, flags ...*FlagSet) {
	// 尝试执行结构体绑定操作，捕获可能发生的错误。
	if err := StructBindE(structPtr, flags...); err != nil {
		// 如果有错误发生，将错误信息打印到标准错误输出，并退出程序。
		fmt.Fprintf(os.Stderr, "flags: %s\n", err)
		os.Exit(1)
	}
}

// StructBindE 通过反射将命令行参数绑定到结构体指针上。
//
// 参数:
//   - structPtr: 任意类型的结构体指针，将以此为模板解析命令行参数。
//   - flags: 可选的FlagSet数组，用于定义命令行参数的解析规则。如果不提供，则使用默认规则。
func StructBindE(structPtr any, flags ...*FlagSet) (err error) {
	// 如果没有提供FlagSet，则使用默认的FlagSet
	if len(flags) == 0 {
		flags = append(flags, Default())
	}

	// 通过反射获取结构体指针的值，并确保可以对其进行设置
	v := rVal(structPtr, true)

	if !v.CanSet() {
		err = fmt.Errorf("cannot set %T", structPtr)
		return
	}

	// 获取结构体所有字段，并准备进行绑定
	var fields []*flagField

	if fields, err = getFields(v, true); err != nil {
		return
	}

	// 遍历所有字段，为每个字段绑定命令行参数
	for _, field := range fields {
		field.UpdateFromEnv() // 从环境变量中更新字段值

		usage := field.Usage // 获取字段的使用说明
		if usage == "" {
			usage = field.Field.Name // 如果没有提供使用说明，则使用字段名
		}

		// 如果字段设置了环境变量，将环境变量信息添加到使用说明中
		if len(field.Env) > 0 {
			usage += fmt.Sprintf(" (env: %s)", strings.Join(field.Env, ", "))
		}

		// 创建并配置命令行参数项
		item := flagSet(flags...).VarPF(field.Value, field.Name, field.Shorthand, usage)
		item.Deprecated = field.Deprecated               // 设置字段的弃用信息
		item.ShorthandDeprecated = field.ShortDeprecated // 设置字段的简写弃用信息

		// 如果字段类型为布尔型，则设置无选项默认值为true
		if fv := reflect.Indirect(field.Value.v); fv.Kind() == reflect.Bool {
			item.NoOptDefVal = "true"
		}
	}

	return
}

// StructPrint 打印结构体指针的字段信息
//
// 参数:
//   - structPtr: 任意类型的结构体指针，函数将通过反射获取其字段信息进行打印
//   - print: 一个函数，用于输出字符串，类似于fmt.Println，用于打印字段信息
func StructPrint(structPtr any, print func(s string)) {
	// 通过反射获取结构体的字段信息
	fields, err := getFields(structPtr)
	if err != nil {
		// 如果获取字段信息过程中出现错误，打印错误信息并返回
		print(err.Error())
		return
	}

	// 计算字段名的最大长度，用于格式化打印
	max := 0
	for _, f := range fields {
		if l := len(f.Field.Name); l > max {
			max = l
		}
	}

	// 格式化并打印每个字段的名称和值
	for _, f := range fields {
		print(fmt.Sprintf("%-*s | %s", max, f.Name, strings.Join(f.Value.Args(), ", ")))
	}
}

// StructToArgs 函数将结构体指针转换为命令行参数字符串切片。
//
// 参数:
//   - structPtr 任意类型的结构体指针，函数将通过反射解析其字段和值。
//
// 返回值 :
//
//   - args 是包含结构体字段及其值的字符串切片，格式为 "--字段名 值"。
func StructToArgs(structPtr any) (args []string) {
	// 通过反射获取结构体指针的实际类型，并提取其字段信息
	fields, _ := getFields(structPtr)
	for _, f := range fields {
		// 遍历每个字段的值，并将其转换为命令行参数格式
		for _, s := range f.Value.Args() {
			args = append(args, "--"+f.Name) // 添加字段名前缀 "--"
			args = append(args, s)           // 添加字段值
		}
	}
	return
}
