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

type album struct {
	title  string
	artist string
	tracks []music
}

type artist struct {
	name   string
	tracks []music
}

type (
	errMsg      struct{ err error }
	musicsMsg   struct{ musics []music }
	albumsMsg   struct{ albums []album }
	artistsMsg  struct{ artists []artist }
	finishedMsg struct{}
)

type progressMsg struct {
	elapsed time.Duration
	total   time.Duration
}

type playingMsg struct {
	music      music
	streamer   beep.StreamSeekCloser
	sampleRate beep.SampleRate
}

// list.Item implementation
func (s music) Title() string       { return s.title }
func (s music) Description() string { return s.artist }
func (s music) FilterValue() string { return s.title }

func (a album) Title() string       { return a.title }
func (a album) Description() string { return a.artist }
func (a album) FilterValue() string { return a.title }

func (a artist) Title() string       { return a.name }
func (a artist) Description() string { return "" }
func (a artist) FilterValue() string { return a.name }

func tickCmd(streamer beep.StreamSeekCloser, sr beep.SampleRate) tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		speaker.Lock()
		elapsed := sr.D(streamer.Position()).Round(time.Second)
		total := sr.D(streamer.Len()).Round(time.Second)
		speaker.Unlock()
		return progressMsg{elapsed, total}
	})
}

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

func fetchAlbums() tea.Msg {
	albumsMap := make(map[string]album)
	homeDir, _ := os.UserHomeDir()
	dir := filepath.Join(homeDir, "Music")

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

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

		// check if no album
		if metadata.Album() == "" {
			return nil
		}

		// create a key for map
		albumKey := metadata.AlbumArtist() + " - " + metadata.Album()

		currentMusic := music{
			title:  metadata.Title(),
			artist: metadata.Artist(),
			path:   path,
		}

		// check if the album already exists in map
		if existingAlbum, ok := albumsMap[albumKey]; ok {
			existingAlbum.tracks = append(existingAlbum.tracks, currentMusic)
			albumsMap[albumKey] = existingAlbum
		} else {
			newAlbum := album{
				title:  metadata.Album(),
				artist: metadata.Artist(),
				tracks: []music{currentMusic},
			}
			albumsMap[albumKey] = newAlbum
		}

		return nil
	})

	// convert the map values into a slice of albums
	var albums []album
	for _, album := range albumsMap {
		albums = append(albums, album)
	}

	return albumsMsg{albums}
}

func fetchArtists() tea.Msg {
	artistsMap := make(map[string]artist)
	homeDir, _ := os.UserHomeDir()
	dir := filepath.Join(homeDir, "Music")

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

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

		// create a key for map
		artistKey := metadata.Artist()

		currentMusic := music{
			title:  metadata.Title(),
			artist: metadata.Artist(),
			path:   path,
		}

		// check if the artist already exists in map
		if existingArtist, ok := artistsMap[artistKey]; ok {
			existingArtist.tracks = append(existingArtist.tracks, currentMusic)
			artistsMap[artistKey] = existingArtist
		} else {
			newArtist := artist{
				name:   metadata.Artist(),
				tracks: []music{currentMusic},
			}
			artistsMap[artistKey] = newArtist
		}

		return nil
	})

	// convert the map values into a slice of artists
	var artists []artist
	for _, artist := range artistsMap {
		artists = append(artists, artist)
	}

	return artistsMsg{artists}
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

	speaker.Clear()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	finishedMsgChan := make(chan tea.Msg, 1)

	go func() {
		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
			finishedMsgChan <- finishedMsg{}
		})))
	}()

	return tea.Batch(
		func() tea.Msg { return playingMsg{music: m, streamer: streamer, sampleRate: format.SampleRate} },
		func() tea.Msg { return <-finishedMsgChan },
		tickCmd(streamer, format.SampleRate),
	)()
}
