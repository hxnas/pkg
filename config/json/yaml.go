package json

import (
	"encoding/json"

	"github.com/hxnas/pkg/config"
)

func Unmarshal(v any) config.ReadFunc {
	return func(data []byte) error { return json.Unmarshal(jcTranslate(data), v) }
}

func init() {
	config.Registry("json", Unmarshal, ".json", ".jsonc")
}
