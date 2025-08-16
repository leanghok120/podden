package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	list        list.Model
	width       int
	height      int
	loaded      bool
	playing     bool
	currPlaying music
}

func main() {
	p := tea.NewProgram(model{loaded: false, playing: false}, tea.WithAltScreen())
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

	if m.loaded {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, screenStyle.Render(m.list.View()))
	}
	if m.playing {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, screenStyle.Render(
			fmt.Sprintf("%s\n\n%s\n", title.Render(m.currPlaying.title), m.currPlaying.artist),
		))
	}
	return "loading music..."
}
