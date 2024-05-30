package sys

import (
	"os"
	"slices"
	"strings"

	"github.com/hxnas/pkg/lod"
)

func NewEnv() *Env {
	return Env{}.Init()
}

type Env struct {
	envMap map[string]string
	keys   []string
}

func (e Env) Init() *Env {
	e.envMap = make(map[string]string)
	return &e
}

func (e *Env) Each(walkFn func(k, v string)) {
	for _, k := range e.keys {
		walkFn(k, e.envMap[k])
	}
}

func (e *Env) Set(k, v string) *Env {
	if k = strings.TrimSpace(k); k != "" {
		idx := slices.Index(e.keys, k)
		if v == "" {
			delete(e.envMap, k)
			if idx > -1 {
				e.keys = slices.Delete(e.keys, idx, idx+1)
			}
		} else {
			e.envMap[k] = v
			if idx == -1 {
				e.keys = append(e.keys, k)
			}
		}
	}

	return e
}

func (e *Env) SetOptional(k, v string) *Env {
	if k = strings.TrimSpace(k); k != "" && v != "" {
		e.envMap[k] = v
		if i := slices.Index(e.keys, k); i == -1 {
			e.keys = append(e.keys, k)
		}
	}
	return e
}

func (e *Env) Append(envs ...string) *Env {
	for _, it := range envs {
		k, v, found := strings.Cut(it, "=")
		if !found {
			k = it
		}
		e.Set(k, v)
	}
	return e
}

func (e *Env) AppendOS() *Env {
	return e.Append(os.Environ()...)
}

func (e *Env) Environ() (environs []string) {
	if e != nil {
		environs = lod.Map(e.keys, func(k string) string { return k + "=" + e.envMap[k] })
	}
	return
}

func (e *Env) Merge(another *Env) {
	another.Each(func(k, v string) { e.Set(k, v) })
}
