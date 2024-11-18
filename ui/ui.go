package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/s992/nfl/shared"
	"github.com/s992/nfl/ui/feed"
	"github.com/s992/nfl/ui/gamepicker"
)

type (
	model struct {
		screen screen
		height int

		gamepicker gamepicker.Model
		feed       feed.Model
		help       help.Model
	}

	screen uint

	keyMap struct {
		Quit key.Binding
	}
)

var (
	defaultKeyMap = keyMap{
		Quit: key.NewBinding(key.WithKeys("q", "ctrl-c"), key.WithHelp("q", "quit")),
	}
)

const (
	screenGamepicker screen = iota
	screenFeed       screen = iota + 1
)

func Run() {
	m := initModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

func initModel() model {
	m := model{}
	m.screen = screenGamepicker
	m.gamepicker = gamepicker.New()
	m.feed = feed.New()
	m.help = help.New()

	return m
}

func (m model) ShortHelp() []key.Binding {
	bindings := []key.Binding{
		defaultKeyMap.Quit,
		shared.SharedKeyMap.Refresh,
	}

	switch m.screen {
	case screenFeed:
		bindings = append(bindings, m.feed.Keybinds()...)
	case screenGamepicker:
		bindings = append(bindings, m.gamepicker.Keybinds()...)
	}

	return bindings
}

func (m model) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.gamepicker.Init(), m.feed.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmds = append(cmds, m.handleSize(msg))
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, defaultKeyMap.Quit):
			return m, tea.Quit
		}
	case gamepicker.WatchGameMsg:
		cmds = append(cmds, m.handleWatch(msg))
	case feed.ExitWatchMsg:
		cmds = append(cmds, m.handleExitWatch())
	}

	var cmd tea.Cmd

	if m.screen == screenGamepicker {
		m.gamepicker, cmd = m.gamepicker.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.screen == screenFeed {
		m.feed, cmd = m.feed.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var subview string

	switch m.screen {
	case screenFeed:
		subview = m.feed.View()
	case screenGamepicker:
		subview = m.gamepicker.View()
	}

	helpAvailableHeight := m.height - lipgloss.Height(subview)
	help := lipgloss.PlaceVertical(helpAvailableHeight, lipgloss.Bottom, m.help.View(m))

	return lipgloss.JoinVertical(lipgloss.Left, subview, help)
}

func (m *model) handleSize(msg tea.WindowSizeMsg) tea.Cmd {
	availableHeight := msg.Height - 1

	m.height = msg.Height
	m.help.Width = msg.Width

	m.gamepicker.SetHeight(availableHeight)
	m.gamepicker.SetWidth(msg.Width)
	m.feed.SetHeight(availableHeight)
	m.feed.SetWidth(msg.Width)

	return nil
}

func (m *model) handleWatch(msg gamepicker.WatchGameMsg) tea.Cmd {
	m.screen = screenFeed

	return m.feed.SetGame(msg.Game)
}

func (m *model) handleExitWatch() tea.Cmd {
	m.screen = screenGamepicker

	return nil
}
