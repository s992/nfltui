package gamepicker

import (
	"fmt"
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/s992/nfl/api"
	"github.com/s992/nfl/shared"
)

type (
	Model struct {
		list    list.Model
		loading bool
		date    time.Time
		err     string
	}

	WatchGameMsg struct {
		Game api.Scoreboard
	}

	refreshGamesMsg struct {
		games []api.Scoreboard
		err   error
	}

	keyMap struct {
		Watch key.Binding
	}
)

var (
	defaultKeyMap = keyMap{
		Watch: key.NewBinding(key.WithKeys("enter", "w"), key.WithHelp("↵/w", "watch")),
	}
)

func New() Model {
	m := Model{}
	m.loading = false
	m.date = time.Now()
	m.list = list.New(nil, list.NewDefaultDelegate(), 0, 0)
	m.list.DisableQuitKeybindings()
	m.list.SetShowHelp(false)
	m.list.SetFilteringEnabled(false)

	return m
}

func (m *Model) SetHeight(height int) {
	m.list.SetHeight(height)
}

func (m *Model) SetWidth(width int) {
	m.list.SetWidth(width)
}

func (m Model) Init() tea.Cmd {
	return m.handleRefreshRequested()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, shared.SharedKeyMap.Refresh):
			cmds = append(cmds, m.handleRefreshRequested())
		case key.Matches(msg, defaultKeyMap.Watch):
			cmds = append(cmds, m.handleWatch())
		}
	case shared.TickMsg:
		cmds = append(cmds, m.handleRefreshRequested())
	case refreshGamesMsg:
		cmds = append(cmds, m.handleFreshGames(msg))
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.list.View()
}

func (m Model) Keybinds() []key.Binding {
	return []key.Binding{
		defaultKeyMap.Watch,
	}
}

func (m *Model) handleFreshGames(next refreshGamesMsg) tea.Cmd {
	if next.err != nil {
		m.list.Title = fmt.Sprintf("Failed to fetch games: %s", next.err)
		return shared.Tick(time.Minute)
	}

	m.list.Title = fmt.Sprintf("NFL : %s", m.date.Format("02 Jan"))
	slices.SortFunc(next.games, func(a, b api.Scoreboard) int {
		if a.State == b.State {
			return 0
		}

		if a.State == api.GameStatePost {
			return 1
		}

		if a.State == api.GameStateInProgress {
			return -1
		}

		if a.State == api.GameStatePre && b.State != api.GameStateInProgress {
			return -1
		}

		return 0
	})

	var items []list.Item

	for _, game := range next.games {
		items = append(items, item{game})
	}

	m.loading = false
	setcmd := m.list.SetItems(items)

	return tea.Batch(setcmd, shared.Tick(time.Minute))
}

func (m *Model) handleRefreshRequested() tea.Cmd {
	if !m.loading {
		m.loading = true

		return func() tea.Msg {
			games, err := api.GetScoreboards(m.date)
			return refreshGamesMsg{games, err}
		}
	}

	return nil
}

func (m Model) handleWatch() tea.Cmd {
	idx := m.list.Index()
	item := m.list.Items()[idx].(item)

	return func() tea.Msg {
		return WatchGameMsg{Game: item.Scoreboard}
	}
}
