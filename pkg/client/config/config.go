package config

import (
	"embed"

	"github.com/knadh/koanf/v2"
)

const (
	ext    = ".toml"
	prefix = "config"
)

//go:embed *.toml
var configFs embed.FS

type config struct {
	configLoader *koanf.Koanf
	fs           *embed.FS
	fname        string
}
