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

// Version函数用于设置或获取版本号。
//
//	如果提供了参数，则尝试更新版本号为第一个非空字符串。
//
// 参数:
//   - updateVer ... - 可变数量的字符串参数，用于更新版本号。第一个非空字符串将被用作新的版本号。
//
// 返回值:
//   - string - 当前的版本号。
func Version(updateVer ...string) string {
	for _, ver := range updateVer {
		// 遍历提供的版本号参数，更新版本号为第一个非空字符串
		if ver != "" {
			version = ver
		}
	}
	return version
}

// ParseFlag 解析命令行参数。
//
// 参数:
//   - flag *FlagSet: 用于解析命令行参数的 FlagSet 实例。
//   - args []string: 命令行参数数组。
func ParseFlag(flag *FlagSet, args []string) (err error) {
	// 初始化命令行工具名称和错误输出目标。
	name, out := name(), os.Stderr

	// 设置当请求帮助时的自定义错误信息。
	pflag.ErrHelp = fmt.Errorf("use %s [...OPTIONS] to start", name)

	// 初始化 flag 配置。
	flag.Init(name, pflag.ContinueOnError)
	flag.SetOutput(out)
	flag.SortFlags = false

	// 设置 flag 使用说明。
	flag.Usage = func() {
		// 输出程序名称和版本信息（如果有）。
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

	// 如果 flag 为默认 FlagSet，则重写 pflag 的 Usage 方法。
	if flag == Default() {
		pflag.Usage = flag.Usage
	}

	// 定义版本信息标志并绑定到 flag。
	const versionFlag = "version"
	if version != "" && flag.Lookup(versionFlag) == nil {
		var shorthand string
		// 为版本信息标志选择一个简写形式，如果未被占用则使用 "v" 或 "V"。
		if flag.Lookup("v") == nil {
			shorthand = "v"
		}
		if shorthand == "" && flag.Lookup("V") == nil {
			shorthand = "V"
		}
		flag.BoolP(versionFlag, shorthand, false, "显示版本信息")
	}

	// 解析命令行参数。
	if err = flag.Parse(args); err != nil {
		return
	}

	// 如果命令行中请求了版本信息，则输出并退出。
	if ver, _ := flag.GetBool(versionFlag); ver {
		fmt.Fprintf(out, "%s", name)
		if version != "" {
			fmt.Fprintf(out, " -- version %s", version)
		}
		os.Exit(0)
	}

	return
}

// Parse 函数解析命令行参数。
//   - 它首先尝试使用默认配置解析命令行参数（从os.Args[1:]开始）。
//   - 如果解析过程中出现错误，将错误信息打印到标准错误输出，并退出程序，退出码为1。
func Parse() {
	// 使用默认配置解析命令行参数
	if err := ParseFlag(Default(), os.Args[1:]); err != nil {
		// 如果解析失败，将错误信息打印到标准错误输出，并退出程序
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
