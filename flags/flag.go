package flags

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

type FlagSet = pflag.FlagSet

var Default = func() *FlagSet { return pflag.CommandLine }

var version string

// 设置或者获取版本号， updateVer参数的第一个不为空的值将设置到版本号中，返回最终版本号
func Version(updateVer ...string) string {
	for _, ver := range updateVer {
		if ver != "" {
			version = ver
		}
	}
	return version
}

func ParseFlag(flag *FlagSet, args []string) (err error) {
	name, out := name(), os.Stderr

	pflag.ErrHelp = fmt.Errorf("use %s [...OPTIONS] to start", name)

	flag.Init(name, pflag.ContinueOnError)
	flag.SetOutput(out)
	flag.SortFlags = false

	flag.Usage = func() {
		fmt.Fprintf(out, "%s", name)
		if version != "" {
			fmt.Fprintf(out, " -- version %s", version)
		}
		fmt.Fprintf(out, "\n\n")
		fmt.Fprintf(out, "USAGE:\n")
		fmt.Fprintf(out, "  %s [...OPTIONS]\n\n", name)
		fmt.Fprintf(out, "OPTIONS:\n")
		fmt.Fprintln(out, flag.FlagUsagesWrapped(0))
	}

	if flag == Default() {
		pflag.Usage = flag.Usage
	}

	const versionFlag = "version"

	if version != "" && flag.Lookup(versionFlag) == nil {
		var shorthand string
		if flag.Lookup("v") == nil {
			shorthand = "v"
		}
		if shorthand == "" && flag.Lookup("V") == nil {
			shorthand = "V"
		}
		flag.BoolP(versionFlag, shorthand, false, "显示版本信息")
	}

	if err = flag.Parse(args); err != nil {
		return
	}

	if ver, _ := flag.GetBool(versionFlag); ver {
		fmt.Fprintf(out, "%s", name)
		if version != "" {
			fmt.Fprintf(out, " -- version %s", version)
		}
		os.Exit(0)
	}

	return
}

func Parse() {
	if err := ParseFlag(Default(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func name() string { return filepath.Base(os.Args[0]) }

func flagSet(flags ...*FlagSet) *FlagSet {
	for _, f := range flags {
		if f != nil {
			return f
		}
	}
	return Default()
}
