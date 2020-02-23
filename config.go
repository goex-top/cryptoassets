package main

import (
	"github.com/BurntSushi/toml"
)

func loadConfig() (Config, error) {
	var conf Config
	_, err := toml.DecodeFile("config.toml", &conf)
	return conf, err
}
