package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/willtrojniak/homegoing/dotmodels"
)

type keyMap struct {
	Help key.Binding
	Quit key.Binding
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit")),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help")),
}

func (m App) ShortHelp() []key.Binding {
	return []key.Binding{m.keys.Quit, m.config.Keys.Up, m.config.Keys.Down, m.config.Keys.Refresh, m.config.Keys.Link, m.config.Keys.Unlink}
}

func (m App) FullHelp() [][]key.Binding {
	return [][]key.Binding{{m.keys.Quit, m.keys.Help}, {m.config.Keys.Refresh}, {m.config.Keys.Link, m.config.Keys.Unlink}}
}

type App struct {
	help   help.Model
	keys   keyMap
	height int

	config dotmodels.DotConfigModel

	error
	isQuitting bool
}

func newApp(configFilePath string) *App {
	help := help.New()
	help.ShowAll = false

	return &App{
		help:   help,
		keys:   keys,
		config: dotmodels.NewDotConfigModel(configFilePath),
	}
}

func (m App) Init() tea.Cmd {
	return m.config.Init()
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case error:
		// TODO: Log other errors
		m.error = msg
		return m, nil
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.help.Width = msg.Width
		// TODO: Handle resizing
	case tea.KeyMsg:
		if m.error != nil {
			m.error = nil
		}
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.isQuitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			// TODO: Implement
		}
	}
	var cmd tea.Cmd
	m.config, cmd = m.config.Update(msg)

	return m, cmd
}

func (m App) View() string {
	var s strings.Builder

	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).AlignHorizontal(lipgloss.Right).Width(m.help.Width)
	versionView := versionStyle.Render("homegoing v0.1.1")
	s.WriteString(versionView)
	s.WriteString("\n")

	if m.isQuitting {
		if m.error != nil {
			s.WriteString(fmt.Sprintf("A fatal error occurred: %v", m.error))
			return s.String()
		}
		return ""
	}

	configView := m.config.View()
	helpView := m.help.View(m)
	versionHeight := strings.Count(versionView, "\n")
	configHeight := strings.Count(configView, "\n")
	helpHeight := strings.Count(helpView, "\n")
	errHeight := 0
	s.WriteString(configView)
	if m.error != nil {
		s.WriteString(fmt.Sprintf("  An error has occurred:\n  %v", m.error))
		errHeight = 1
	}
	s.WriteString(strings.Repeat("\n", max(m.height-configHeight-versionHeight-helpHeight-errHeight-2, 0)))
	s.WriteString(helpView)
	return s.String()
}
