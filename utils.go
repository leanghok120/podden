package main

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type lyricLine struct {
	Time float64
	Text string
}

// styles
var (
	screenStyle          lipgloss.Style
	titleStyle           lipgloss.Style
	titleBackgroundStyle lipgloss.Style
	artistStyle          lipgloss.Style
	lyricStyle           lipgloss.Style
	timeStyle            lipgloss.Style
)

func initStyles() {
	screenStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(30).
		MaxWidth(35).
		Height(14)

	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(cfg.TitleForeground)).
		Background(lipgloss.Color("62")).
		Padding(0, 1)

	artistStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")) // muted gray

	lyricStyle = lipgloss.NewStyle().
		Width(26).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("252")).
		Italic(true)

	timeStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
}

// helper functions
// place content in the center
func (m model) center(content string) string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

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

// handle album selection in list
func (m model) handleAlbumSelection() model {
	if selected, ok := m.list.SelectedItem().(album); ok {
		items := make([]list.Item, len(selected.tracks))
		for i, track := range selected.tracks {
			items[i] = track
		}
		m.list.SetItems(items)
		m.list.Title = selected.title
		m.showAlbums = false
		m.loaded = true
		m.list.SetFilterState(list.Unfiltered)
	}
	return m
}

// handle artist selection in list
func (m model) handleArtistSelection() model {
	if selected, ok := m.list.SelectedItem().(artist); ok {
		items := make([]list.Item, len(selected.tracks))
		for i, track := range selected.tracks {
			items[i] = track
		}
		m.list.SetItems(items)
		m.list.Title = selected.name
		m.showArtists = false
		m.loaded = true
		m.list.SetFilterState(list.Unfiltered)
	}
	return m
}

func parseLRC(raw string) ([]lyricLine, error) {
	var lyrics []lyricLine

	// regex to match [mm:ss.xx]
	re := regexp.MustCompile(`\[(\d+):(\d+\.\d+)\](.*)`)

	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 4 {
			minutes, _ := strconv.Atoi(matches[1])
			seconds, _ := strconv.ParseFloat(matches[2], 64)
			text := strings.TrimSpace(matches[3])
			totalTime := float64(minutes)*60 + seconds
			lyrics = append(lyrics, lyricLine{Time: totalTime, Text: text})
		}
	}

	return lyrics, nil
}
