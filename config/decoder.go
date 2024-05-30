package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var decoder = DecoderFactory{}

func Decode(source string, v any) error { return decoder.Decode(source, v) }

func Registry(name string, unmarshal DecodeFunc, exts ...string) {
	decoder.Register(name, unmarshal, exts...)
}

type (
	DecoderFactory map[string]DecodeFunc
	ReadFunc       = func(data []byte) (err error)
	DecodeFunc     = func(value any) ReadFunc
)

func (f DecoderFactory) Register(name string, decodeFunc DecodeFunc, exts ...string) {
	name = strings.ToLower(name)
	f[name] = decodeFunc
	for _, ext := range exts {
		f[strings.ToLower(ext)] = decodeFunc
	}
}

func (f DecoderFactory) Decode(source string, v any) (err error) {
	if source == "" {
		err = fmt.Errorf("source is empty")
		return
	}

	var unmarshal DecodeFunc
	var path string

	if n, p, ok := strings.Cut(source, ":"); ok {
		if unmarshal = f[n]; unmarshal != nil {
			path = p
		}
	}

	if unmarshal == nil {
		if unmarshal = f[filepath.Ext(source)]; unmarshal != nil {
			path = source
		}
	}

	if unmarshal == nil {
		err = fmt.Errorf("unsupport source: %s", source)
		return
	}

	return f.readBytes(path, unmarshal(v))
}

func (f DecoderFactory) readBytes(path string, read ReadFunc) (err error) {
	if path == "" {
		err = fmt.Errorf("source path is empty")
		return
	}

	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return
	}

	return read(data)
}
