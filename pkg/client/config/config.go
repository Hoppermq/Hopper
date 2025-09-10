// Package config represent the configuration package.
package config

import (
	"embed"
	"fmt"

	"github.com/zixyos/goloader/config"
)


//go:embed *.toml
var configFs embed.FS
type ClientConfig struct {
	Environment string `koanf:"environment"`
	Name        string `koanf:"name"`
	Version     string `koanf:"version"`
	MaxRetries  int    `koanf:"maxRetries"`
	Timeout     int    `koanf:"timeout"`
}

type TransportConfig struct {
	Protocol  string `koanf:"protocol"`
	Host      string `koanf:"host"`
	Port      int    `koanf:"port"`
	Keepalive string `koanf:"keepalive"`
	Heartbeat string `koanf:"heartbeat"`
}

type SecurityConfig struct{}

type AuthConfig struct{}

type ClientConfiguration struct {
	Client    ClientConfig    `koanf:"client"`
	Transport TransportConfig `koanf:"transport"`
	Auth      AuthConfig      `koanf:"auth"`
	Security  SecurityConfig  `koanf:"security"`
}

type Option func(*ClientConfiguration)

func (c *ClientConfiguration) PostLoad() error {
	return nil
}

func LoadConfig() (*ClientConfiguration, error) {
	var cfg ClientConfiguration

	if err := config.Load(&cfg, config.WithFs(configFs)); err != nil {
		fmt.Printf("error loading client config: %v\n", err)
		return nil, err
	}

	return &cfg, nil
}
