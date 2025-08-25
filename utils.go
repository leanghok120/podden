package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/0xAX/notificator"
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
	helpMenu             lipgloss.Style
)

func fallbackColor(value, def string) lipgloss.Color {
	if value == "" {
		return lipgloss.Color(def)
	}
	return lipgloss.Color(value)
}

// returns either lipgloss.Color or lipgloss.AdaptiveColor
func fallbackAdaptiveColor(value string, def lipgloss.AdaptiveColor) lipgloss.TerminalColor {
	if value == "" {
		return def
	}
	return lipgloss.Color(value)
}

func initStyles() {
	screenStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(cfg.BorderForeground)).
		Padding(1, 2).
		Width(30).
		MaxWidth(35).
		Height(14)

	titleStyle = lipgloss.NewStyle().
		Foreground(fallbackColor(cfg.HeadingForeground, "230")).
		Background(fallbackColor(cfg.HeadingBackground, "62")).
		Padding(0, 1)

	artistStyle = lipgloss.NewStyle().
		Foreground(fallbackColor(cfg.ArtistForeground, "243")) // muted gray

	lyricStyle = lipgloss.NewStyle().
		Width(26).
		Align(lipgloss.Center).
		Foreground(fallbackColor(cfg.LyricsForeground, "252")).
		Italic(true)

	timeStyle = lipgloss.NewStyle().
		Foreground(fallbackColor(cfg.TimeForeground, "240"))

	helpMenu = lipgloss.NewStyle().
		Padding(0, 1)
}

// helper functions
// update the styles of bubbles component
func setCustomBubblesStyle() list.Styles {
	styles := list.DefaultStyles()

	styles.Title = lipgloss.NewStyle().
		Background(fallbackColor(cfg.HeadingBackground, "62")).
		Foreground(fallbackColor(cfg.HeadingForeground, "230")).
		Padding(0, 1)

	return styles
}

func customDelegate() list.ItemDelegate {
	delegate := list.NewDefaultDelegate()
	s := &delegate.Styles

	s.NormalTitle = lipgloss.NewStyle().
		Foreground(fallbackAdaptiveColor(cfg.NormalTitleForeground,
					lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"})).
		Padding(0, 0, 0, 2) //nolint:mnd

	s.NormalDesc = s.NormalTitle.
		Foreground(fallbackAdaptiveColor(cfg.NormalDescForeground,
			lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}))

	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(fallbackAdaptiveColor(cfg.SelectedTitleBorderForeground,
			lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"})).
		Foreground(fallbackAdaptiveColor(cfg.SelectedTitleForeground,
			lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"})).
		Padding(0, 0, 0, 1)

	s.SelectedDesc = s.SelectedTitle.
		Foreground(fallbackAdaptiveColor(cfg.SelectedDescForeground,
			lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}))

	s.DimmedTitle = lipgloss.NewStyle().
		Foreground(fallbackAdaptiveColor(cfg.DimmedTitleForeground,
					lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})).
		Padding(0, 0, 0, 2) //nolint:mnd

	s.DimmedDesc = s.DimmedTitle.
		Foreground(fallbackAdaptiveColor(cfg.DimmedDescForeground,
			lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"}))

	return delegate
}

// place content in the center and add a help menu
func (m model) center(content string) string {
	screen := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)

	if cfg.ShowHelp {
		return lipgloss.JoinVertical(lipgloss.Left, screen, helpMenu.Render(m.help.View(keys)))
	}
	return screen
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

func sendNotification(m music, body string) {
	f, err := os.CreateTemp("", "cover-image")
	if err != nil {
		return
	}
	defer os.Remove(f.Name())

	_, err = f.Write(m.cover)
	if err != nil {
		return
	}
	f.Close()

	notifyTitle := fmt.Sprintf("%s - %s", m.title, m.artist)
	notify.Push(notifyTitle, body, f.Name(), notificator.UR_NORMAL)
}
