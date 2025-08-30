package config

import (
	"embed"
)

// Option type is a function that modify the configuration.
type Option func(*config) error

// WithFs set the embed fs to the config struct.
func WithFs(fs embed.FS) Option {
	return func(c *config) error {
		c.fs = &fs
		return nil
	}
}

// WithFName function set the file name to the config struct.
func WithFName(fname string) Option {
	return func(c *config) error {
		c.fname = fname
		return nil
	}
}
