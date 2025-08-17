package main

import (
	"fmt"

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
	showAlbums  bool
	playing     bool
	paused      bool
	currPlaying music
	streamer    beep.StreamSeekCloser
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

			case "enter":
				// handle album selection
				if m.showAlbums {
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
						return m, nil
					}
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
				m.showAlbums = false

				return m, func() tea.Msg { return fetchMusics() }

			case "a":
				m.playing = false
				m.loaded = false

				return m, func() tea.Msg { return fetchAlbums() }
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

	case albumsMsg:
		items := make([]list.Item, len(msg.albums))
		for i, m := range msg.albums {
			items[i] = m
		}
		l := list.New(items, list.NewDefaultDelegate(), 30, 10)
		l.Title = "Albums"

		m.list = l
		m.showAlbums = true

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

	if m.loaded || m.showAlbums {
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
	currList := m.list.View()
	playing := fmt.Sprintf("%s\n\n%s\n", titleStyle.Render(m.currPlaying.title), m.currPlaying.artist)

	if m.loaded {
		return m.center(screenStyle.Render(currList))
	}
	if m.playing {
		return m.center(screenStyle.Render(playing))
	}
	if m.showAlbums {
		return m.center(screenStyle.Render(currList))
	}
	return m.center(screenStyle.Render("loading music..."))
}
