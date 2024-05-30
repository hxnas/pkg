package config

import (
	"os"
)

type ConfigFile string

type ConfigFileValue struct {
	path    string
	defPath string
	referer any
}

func (b *ConfigFileValue) String() string { return b.defPath }
func (b *ConfigFileValue) Type() string   { return "configfile" }
func (b *ConfigFileValue) Set(s string) (err error) {
	if b.path = s; b.path != "" {
		if err = Decode(s, b.referer); os.IsNotExist(err) {
			err = nil
		}
	}
	return
}
