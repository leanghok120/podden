package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// styles
var (
	screenStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Width(30).
			MaxWidth(35).
			Height(14)

	titleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1)

	// bratStyle := lipgloss.NewStyle().
	// 	Padding(1, 2).
	// 	Width(32).
	// 	MaxWidth(35).
	// 	Height(16).
	// 	Background(lipgloss.Color("#ffffff"))

)

// place content in the center
func (m model) center(content string) string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

// helper functions
// play next song
func (m model) nextSong(l list.Model) (list.Model, tea.Cmd) {
	l.CursorDown()
	selected, ok := l.SelectedItem().(music)
	if !ok {
		return l, nil
	}
	return l, func() tea.Msg { return playMusic(selected) }
}

// play previous song
func (m model) prevSong(l list.Model) (list.Model, tea.Cmd) {
	l.CursorUp()
	selected, ok := l.SelectedItem().(music)
	if !ok {
		return l, nil
	}
	return l, func() tea.Msg { return playMusic(selected) }
}
