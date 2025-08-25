package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/0xAX/notificator"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	musicDirFlag = flag.String("m", "", "set your music directory (the directory where all your musics are in)")
	cfg          config
	notify       *notificator.Notificator
)

func main() {
	flag.Parse()
	notify = notificator.New(notificator.Options{})
	loadConfig(&cfg)
	initStyles()

	p := tea.NewProgram(initModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
