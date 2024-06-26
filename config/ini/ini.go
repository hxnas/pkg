package ini

import (
	"github.com/hxnas/pkg/config"
	"gopkg.in/ini.v1"
)

func Unmarshal(v any) config.ReadFunc {
	return func(data []byte) error { return ini.MapTo(v, data) }
}

func init() {
	config.Registry("ini", Unmarshal, ".ini")
}
