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
	BorderForeground              string `yaml:"border_foreground"`
	NormalTitleForeground         string `yaml:"normal_title_foreground"`
	NormalDescForeground          string `yaml:"normal_desc_foreground"`
	SelectedTitleBorderForeground string `yaml:"selected_title_border_foreground"`
	SelectedTitleForeground       string `yaml:"selected_title_foreground"`
	SelectedDescForeground        string `yaml:"selected_desc_foreground"`
	DimmedTitleForeground         string `yaml:"dimmed_title_foreground"`
	DimmedDescForeground          string `yaml:"dimmed_desc_foreground"`
	ArtistForeground              string `yaml:"artist_foreground"`
	TimeForeground                string `yaml:"time_foreground"`
	LyricsForeground              string `yaml:"lyrics_foreground"`
	LyricsHighlightForeground     string `yaml:"lyrics_highlight_foreground"`
	LyricsHighlightBackground     string `yaml:"lyrics_highlight_background"`
	ShowHelp                      bool   `yaml:"show_help"`
}

var defaultConfigYaml = `# heading styles (album, songs, artists)
heading_background: ""
heading_foreground: ""
border_foreground: ""

# list styles
normal_title_foreground: ""
normal_desc_foreground: ""

selected_title_border_foreground: ""
selected_title_foreground: ""
selected_desc_foreground: ""

dimmed_title_foreground: ""
dimmed_desc_foreground: ""

# playing styles
artist_foreground: ""
time_foreground: ""
lyrics_foreground: ""
lyrics_highlight_foreground: ""
lyrics_highlight_background: ""

show_help: true
`

func loadConfig(cfg *config) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	configPath := filepath.Join(configDir, "podden", "config.yml")

	// check if ~/.config/podden exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(configPath), 0755)
		err = os.WriteFile(configPath, []byte(defaultConfigYaml), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	f, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
