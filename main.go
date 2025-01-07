package main

import (
	"github.com/hoppermq/hopper/internal/config"
)

const appName = "Hopper";

func main() {
  _, err := config.New(appName);
  if err != nil {
    panic(err)
  }
}
