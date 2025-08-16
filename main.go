package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

// helper functions
// place content in the center
func (m model) center(content string) string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

// play next song
func (m model) nextSong(l list.Model) (list.Model, tea.Cmd) {
	// stop currPlaying music
	if m.playing {
		speaker.Clear()
	}

	l.CursorDown()
	selected, ok := l.SelectedItem().(music)
	if !ok {
		return l, nil
	}
	return l, func() tea.Msg { return playMusic(selected) }
}

// play previous song
func (m model) prevSong(l list.Model) (list.Model, tea.Cmd) {
	// stop currPlaying music
	if m.playing {
		speaker.Clear()
	}

	l.CursorUp()
	selected, ok := l.SelectedItem().(music)
	if !ok {
		return l, nil
	}
	return l, func() tea.Msg { return playMusic(selected) }
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
		m.loaded = true
		m.playing = false
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
	// styles
	screenStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(30).
		MaxWidth(35).
		Height(14)

	// bratStyle := lipgloss.NewStyle().
	// 	Padding(1, 2).
	// 	Width(32).
	// 	MaxWidth(35).
	// 	Height(16).
	// 	Background(lipgloss.Color("#ffffff"))

	title := lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	// strings
	songs := m.list.View()
	playing := fmt.Sprintf("%s\n\n%s\n", title.Render(m.currPlaying.title), m.currPlaying.artist)

	if m.loaded {
		return m.center(screenStyle.Render(songs))
	}
	if m.playing {
		return m.center(screenStyle.Render(playing))
	}
	return "loading music..."
}
