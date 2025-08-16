package main

import (
	"io/fs"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhowden/tag"
)

type music struct {
	title  string
	artist string
}

type (
	errMsg    struct{ err error }
	musicsMsg struct{ musics []music }
)

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

		musics = append(musics, music{title: metadata.Title(), artist: metadata.Artist()})
		return nil
	})

	return musicsMsg{musics}
}
