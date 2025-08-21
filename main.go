package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	musicDirFlag = flag.String("m", "", "set your music directory (the directory where all your musics are in)")
	cfg          config
)

func main() {
	flag.Parse()
	loadConfig(&cfg)
	initStyles()

	p := tea.NewProgram(model{loaded: false, playing: false, paused: false}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
