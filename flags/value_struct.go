package flags

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

func ParseStruct(structPtr any) {
	if err := ParseStructE(Default(), structPtr, os.Args[1:]...); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func ParseStructE(flag *FlagSet, structPtr any, args ...string) (err error) {
	if err = StructBindE(structPtr); err != nil {
		return
	}

	if err = ParseFlag(flagSet(flag), args); err != nil {
		return
	}

	return
}

func StructBind(structPtr any, flags ...*FlagSet) {
	if err := StructBindE(structPtr, flags...); err != nil {
		fmt.Fprintf(os.Stderr, "flags: %s\n", err)
		os.Exit(1)
	}
}

func StructBindE(structPtr any, flags ...*FlagSet) (err error) {
	if len(flags) == 0 {
		flags = append(flags, Default())
	}

	v := reflect.Indirect(reflect.ValueOf(structPtr))

	if !v.CanSet() {
		err = fmt.Errorf("cannot set %T", structPtr)
		return
	}

	var fields []*flagField

	if fields, err = getFields(v, true); err != nil {
		return
	}

	for _, field := range fields {
		field.UpdateFromEnv()

		usage := field.Usage
		if usage == "" {
			usage = field.Field.Name
		}

		if len(field.Env) > 0 {
			usage += fmt.Sprintf(" (env: %s)", strings.Join(field.Env, ", "))
		}

		item := flagSet(flags...).VarPF(field.Value, field.Name, field.Shorthand, usage)
		item.Deprecated = field.Deprecated
		item.ShorthandDeprecated = field.ShortDeprecated

		if fv := reflect.Indirect(field.Value.v); fv.Kind() == reflect.Bool {
			item.NoOptDefVal = "true"
		}
	}

	return
}

func StructPrint(structPtr any, print func(s string)) {
	fields, err := getFields(reflect.Indirect(reflect.ValueOf(structPtr)))
	if err != nil {
		print(err.Error())
		return
	}

	max := 0
	for _, f := range fields {
		if l := len(f.Field.Name); l > max {
			max = l
		}
	}

	for _, f := range fields {
		print(fmt.Sprintf("%-*s | %s", max, f.Name, strings.Join(f.Value.Args(), ", ")))
	}
}

func StructToArgs(structPtr any) (args []string) {
	fields, _ := getFields(reflect.Indirect(reflect.ValueOf(structPtr)))
	for _, f := range fields {
		for _, s := range f.Value.Args() {
			args = append(args, "--"+f.Name)
			args = append(args, s)
		}
	}
	return
}
