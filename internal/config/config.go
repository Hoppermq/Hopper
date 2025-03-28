// Package config represent the config package.
package config

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/fs"
	"github.com/knadh/koanf/v2"
)

const (
	envName = "APP_ENV"
	ext     = ".toml"
	prefix  = "config"
)

// PostLoad is for postloafing ahaha.
type PostLoad interface {
	PostLoad() error
}

//go:embed *.toml
var configFs embed.FS

type config struct {
	target       interface{}
	configLoader *koanf.Koanf
	fs           *embed.FS
	fname        string
}

func loader(target interface{}, options ...Option) error {
	conf := config{
		target: target,
	}

	for _, opt := range options {
		if err := opt(&conf); err != nil {
			return err
		}
	}

	return conf.load()
}

func (c *config) loadFromFS() error {
	if err := c.configLoader.Load(
		fs.Provider(*c.fs, c.fname),
		toml.Parser(),
	); err != nil {
		return fmt.Errorf("error ;oading config file from fs: %w", err)
	}
	return nil
}

func (c *config) loadFromLocal() error {
	if err := c.configLoader.Load(
		file.Provider(c.fname),
		toml.Parser(),
	); err != nil {
		return fmt.Errorf("error loading config file from file: %w", err)
	}
	return nil
}

func (c *config) loadFile() error {
	if c.fname == "" {
		return nil
	}

	if c.fs != nil {
		return c.loadFromFS()
	}

	return c.loadFromLocal()
}

func (c *config) load() error {
	c.configLoader = koanf.New(".")

	if err := c.loadFile(); err != nil {
		return fmt.Errorf("error loading file: %w", err)
	}

	if err := c.configLoader.Load(env.Provider("", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(s), "_", ".")
	}), nil); err != nil {
		return fmt.Errorf("error loading env variables: %w", err)
	}

	if err := c.configLoader.Unmarshal("", c.target); err != nil {
		return fmt.Errorf("error while unmarshalling env vars: %w", err)
	}

	if postLoad, ok := c.target.(PostLoad); ok {
		if err := postLoad.PostLoad(); err != nil {
			return fmt.Errorf("error while running post load: %w", err)
		}
	}
	return nil
}

// Configuration represent the application configuration.
type Configuration struct {
	Test string `koanf:"test"`
}

// New create a new configuration from files and env.
func New(_ string) (*Configuration, error) {
	var conf Configuration
	if err := loader(&conf,
		WithFs(configFs),
		WithFName("config."+os.Getenv(envName)+ext),
	); err != nil {
    return nil, fmt.Errorf("failed to load configuration: %w", err)	
  }

	return &conf, nil
}
