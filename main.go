package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type songItem struct {
	title  string
	artist string
}

func (s songItem) Title() string       { return s.title }
func (s songItem) Description() string { return s.artist }
func (s songItem) FilterValue() string { return s.title }

type model struct {
	list   list.Model
	width  int
	height int
}

func main() {
	items := []list.Item{
		songItem{"Song 1", "Artist A"},
		songItem{"Song 2", "Artist B"},
		songItem{"Song 3", "Artist C"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 30, 10)
	l.Title = "Songs"

	m := model{list: l}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
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
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	m.list.SetShowHelp(false)
	return m, cmd
}

func (m model) View() string {
	// styles
	screenStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(30)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, screenStyle.Render(m.list.View()))
}
