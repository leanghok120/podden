package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
)

type model struct {
	list        list.Model
	width       int
	height      int
	loaded      bool
	playing     bool
	paused      bool
	currPlaying music
	streamer    beep.StreamSeekCloser
}

func main() {
	p := tea.NewProgram(model{loaded: false, playing: false, paused: false}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return fetchMusics
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			if selected, ok := m.list.SelectedItem().(music); ok {
				return m, func() tea.Msg { return playMusic(selected) }
			}

		case " ":
			if m.playing {
				if m.paused {
					speaker.Unlock()
					m.paused = false
				} else {
					speaker.Lock()
					m.paused = true
				}
			}

		case "n":
			var cmd tea.Cmd
			m.list, cmd = m.nextSong(m.list)
			return m, cmd

		case "p":
			var cmd tea.Cmd
			m.list, cmd = m.prevSong(m.list)
			return m, cmd

		case "s":
			m.playing = false
			m.loaded = true
		}

	case musicsMsg:
		items := make([]list.Item, len(msg.musics))
		for i, m := range msg.musics {
			items[i] = m
		}
		l := list.New(items, list.NewDefaultDelegate(), 30, 10)
		l.Title = "Songs"

		m.list = l
		m.loaded = true

	case playingMsg:
		m.loaded = false
		m.playing = true
		m.currPlaying = msg.music
		m.streamer = msg.streamer

	case finishedMsg:
		var cmd tea.Cmd
		m.list, cmd = m.nextSong(m.list)
		return m, cmd
	}

	if m.loaded {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		m.list.SetShowHelp(false)
		m.list.SetShowStatusBar(false)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	// strings
	songs := m.list.View()
	playing := fmt.Sprintf("%s\n\n%s\n", titleStyle.Render(m.currPlaying.title), m.currPlaying.artist)

	if m.loaded {
		return m.center(screenStyle.Render(songs))
	}
	if m.playing {
		return m.center(screenStyle.Render(playing))
	}
	return m.center(screenStyle.Render("loading music..."))
}
