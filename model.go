package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
)

type model struct {
	help        help.Model
	list        list.Model
	width       int
	height      int
	loaded      bool
	showAlbums  bool // for albums page
	showArtists bool // for artists page
	playing     bool
	paused      bool
	lyrics      []lyricLine
	currLyric   string
	elapsed     time.Duration
	total       time.Duration
	currPlaying music
	streamer    beep.StreamSeekCloser
	volume      *effects.Volume
	sampleRate  beep.SampleRate
}

func initModel() model {
	help := help.New()

	return model{loaded: false, playing: false, paused: false, help: help}
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
		if m.list.FilterState() != list.Filtering {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit

			case "?":
				m.help.ShowAll = !m.help.ShowAll
				return m, nil

			case "enter":
				// handle album selection
				if m.showAlbums {
					m = m.handleAlbumSelection()
					return m, nil
				}
				// handle artist selection
				if m.showArtists {
					m = m.handleArtistSelection()
					return m, nil
				}
				// handle song selection and playback
				if selected, ok := m.list.SelectedItem().(music); ok {
					return m, func() tea.Msg { return playMusic(selected) }
				}

			case " ":
				if m.paused {
					speaker.Unlock()
					m.paused = false
				} else {
					speaker.Lock()
					m.paused = true
					sendNotification(m.currPlaying, "paused")
				}

			case ">", "right":
				samplesToJump := m.sampleRate.N(5 * time.Second)
				m.streamer.Seek(m.sampleRate.N(m.elapsed) + samplesToJump)

			case "<", "left":
				samplesToJump := m.sampleRate.N(5 * time.Second)
				newPos := m.sampleRate.N(m.elapsed) - samplesToJump
				if newPos < 0 {
					newPos = 0
				}
				m.streamer.Seek(newPos)

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
				m.showAlbums = false
				m.showArtists = false

				return m, func() tea.Msg { return fetchMusics() }

			case "a":
				m.playing = false
				m.loaded = false
				m.showArtists = false

				return m, func() tea.Msg { return fetchAlbums() }

			case "d":
				m.playing = false
				m.loaded = false
				m.showAlbums = false

				return m, func() tea.Msg { return fetchArtists() }

			case "f":
				m.playing = true
				m.loaded = false
				m.showArtists = false
				m.showAlbums = false

			case "+":
				speaker.Lock()
				m.volume.Volume += 0.5
				speaker.Unlock()

			case "-":
				speaker.Lock()
				m.volume.Volume -= 0.5
				speaker.Unlock()
			}
		}

	case musicsMsg:
		items := make([]list.Item, len(msg.musics))
		for i, m := range msg.musics {
			items[i] = m
		}
		l := list.New(items, customDelegate(), 30, 10)
		l.Title = "Songs"
		l.Styles = setCustomBubblesStyle()

		m.list = l
		m.loaded = true

	case albumsMsg:
		items := make([]list.Item, len(msg.albums))
		for i, a := range msg.albums {
			items[i] = a
		}
		l := list.New(items, customDelegate(), 30, 10)
		l.Title = "Albums"
		l.Styles = setCustomBubblesStyle()

		m.list = l
		m.showAlbums = true

	case artistsMsg:
		items := make([]list.Item, len(msg.artists))
		for i, a := range msg.artists {
			items[i] = a
		}
		l := list.New(items, customDelegate(), 30, 10)
		l.Title = "Artists"
		l.Styles = setCustomBubblesStyle()

		m.list = l
		m.showArtists = true

	case playingMsg:
		m.loaded = false
		m.playing = true
		m.currPlaying = msg.music
		m.streamer = msg.streamer
		m.volume = msg.volume
		m.sampleRate = msg.sampleRate
		m.lyrics = nil // Reset lyrics for the new song
		m.currLyric = "♪"
		m.paused = false
		m.elapsed = 0
		m.total = 0

	case progressMsg:
		m.elapsed = msg.elapsed
		m.total = msg.total

		for _, l := range m.lyrics {
			if m.elapsed.Seconds() >= l.Time {
				m.currLyric = l.Text
			} else {
				break
			}
		}

		return m, tickCmd(m.streamer, m.sampleRate)

	case lyricsMsg:
		m.lyrics = msg.lyrics
		if len(m.lyrics) > 0 && m.lyrics[0].Time > 0 {
			m.currLyric = "♪"
		}

	case finishedMsg:
		var cmd tea.Cmd
		m.list, cmd = m.nextSong(m.list)
		return m, cmd
	}

	if m.loaded || m.showAlbums || m.showArtists {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		m.list.SetShowHelp(false)
		m.list.SetShowStatusBar(false)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	if m.playing {
		title := titleStyle.Render(m.currPlaying.title)
		artist := artistStyle.Render(m.currPlaying.artist)

		timeInfo := timeStyle.Render(fmt.Sprintf("%s / %s", m.elapsed, m.total))
		lyric := lyricStyle.Render(m.currLyric)

		playingContent := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			artist,
			"",
			timeInfo,
			"",
			"",
			lyric,
		)

		return m.center(screenStyle.Render(playingContent))
	}

	if m.loaded || m.showAlbums || m.showArtists {
		return m.center(screenStyle.Render(m.list.View()))
	}

	return m.center(screenStyle.Render("loading..."))
}
