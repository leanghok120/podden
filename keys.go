package main

import "github.com/charmbracelet/bubbles/key"

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Next, k.Prev},
		{k.Albums, k.Songs, k.Artists, k.Playing},
		{k.Play, k.Pause, k.Forward, k.Rewind},
		{k.Help, k.Quit},
	}
}

type keyMap struct {
	// list navigation
	Up   key.Binding
	Down key.Binding

	// playback control
	Play    key.Binding
	Pause   key.Binding
	Next    key.Binding
	Prev    key.Binding
	Forward key.Binding
	Rewind  key.Binding

	// page navigation
	Albums  key.Binding
	Songs   key.Binding
	Artists key.Binding
	Playing key.Binding
	Help    key.Binding
	Quit    key.Binding
}

var keys = keyMap{
	// list navigation
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),

	// playback control
	Play: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "play"),
	),
	Pause: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("space", "pause/resume"),
	),
	Next: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "next song"),
	),
	Prev: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "prev song"),
	),
	Forward: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "fast forward"),
	),
	Rewind: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "rewind"),
	),

	// page navigation
	Albums: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "album"),
	),
	Songs: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "songs"),
	),
	Artists: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "artists"),
	),
	Playing: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "playing"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
}
