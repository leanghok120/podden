package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhowden/tag"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

type music struct {
	title  string
	artist string
	path   string
}

type (
	errMsg      struct{ err error }
	musicsMsg   struct{ musics []music }
	finishedMsg struct{}
)

type playingMsg struct {
	music    music
	streamer beep.StreamSeekCloser
}

// list.Item implementation
func (s music) Title() string       { return s.title }
func (s music) Description() string { return s.artist }
func (s music) FilterValue() string { return s.title }

func fetchMusics() tea.Msg {
	var musics []music
	homeDir, _ := os.UserHomeDir()
	dir := filepath.Join(homeDir, "Music")

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// get common audio file
		ext := filepath.Ext(path)
		if ext != ".mp3" && ext != ".flac" && ext != ".m4a" {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer f.Close()

		metadata, err := tag.ReadFrom(f)
		if err != nil {
			return nil
		}

		musics = append(musics, music{title: metadata.Title(), artist: metadata.Artist(), path: path})
		return nil
	})

	return musicsMsg{musics}
}

func playMusic(m music) tea.Msg {
	f, err := os.Open(m.path)
	if err != nil {
		return errMsg{err}
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		f.Close()
		return errMsg{err}
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	finishedMsgChan := make(chan tea.Msg, 1)

	go func() {
		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
			finishedMsgChan <- finishedMsg{}
		})))
	}()

	return tea.Batch(
		func() tea.Msg { return playingMsg{music: m, streamer: streamer} },
		func() tea.Msg { return <-finishedMsgChan },
	)()
}
