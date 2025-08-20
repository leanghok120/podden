package main

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type config struct {
	TitleForeground string `yaml:"title_foreground"`
}

func loadConfig(cfg *config) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(filepath.Join(configDir, "podden", "config.yml"))
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
