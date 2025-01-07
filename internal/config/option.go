package config

import (
	"embed"
)

// type Option is a function that modify the configuration
type Option func(*config) error;

func WithFs(fs embed.FS) Option {
  return func(c *config) error {
    c.fs = &fs;
    return nil
  }
}

func WithFName(fname string) Option {
  return func(c *config) error {
    c.fname = fname;
    return nil
  }
}
