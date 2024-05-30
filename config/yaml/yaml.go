package yaml

import (
	"github.com/hxnas/pkg/config"
	"gopkg.in/yaml.v3"
)

func Unmarshal(v any) config.ReadFunc {
	return func(data []byte) error { return yaml.Unmarshal(data, v) }
}

func init() {
	config.Registry("yaml", Unmarshal, ".yaml", ".yml")
}
