package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type config struct {
	HeadingForeground             string `yaml:"heading_foreground"`
	HeadingBackground             string `yaml:"heading_background"`
	NormalTitleForeground         string `yaml:"normal_title_foreground"`
	NormalDescForeground          string `yaml:"normal_desc_foreground"`
	SelectedTitleBorderForeground string `yaml:"selected_title_border_foreground"`
	SelectedTitleForeground       string `yaml:"selected_title_foreground"`
	SelectedDescForeground        string `yaml:"selected_desc_foreground"`
	DimmedTitleForeground         string `yaml:"dimmed_title_foreground"`
	DimmedDescForeground          string `yaml:"dimmed_desc_foreground"`
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
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
