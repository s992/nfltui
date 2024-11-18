package feed

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/s992/nfl/api"
	"github.com/s992/nfl/shared"
)

type (
	Model struct {
		game           api.Scoreboard
		plays          []api.Play
		hasInitialData bool
		loading        bool
		height         int
		width          int
	}

	ExitWatchMsg struct{}

	keyMap struct {
		Exit key.Binding
	}

	refreshFeedMsg struct {
		plays []api.Play
		err   error
	}
)

const (
	refreshInterval = time.Minute * 10
)

var (
	defaultKeyMap = keyMap{
		Exit: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back to game list")),
	}

	frameStyle           = lipgloss.NewStyle()
	titleContainerStyle  = lipgloss.NewStyle().Padding(1, 0)
	titleStyle           = lipgloss.NewStyle().AlignHorizontal(lipgloss.Center)
	stampStyle           = lipgloss.NewStyle().PaddingRight(1)
	playDescriptionStyle = lipgloss.NewStyle().PaddingRight(1)
)

func New() Model {
	return Model{}
}

func (m *Model) SetGame(game api.Scoreboard) tea.Cmd {
	m.game = game

	return m.handleRefreshRequested()
}

func (m *Model) SetHeight(height int) {
	m.height = height
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m Model) Keybinds() []key.Binding {
	return []key.Binding{
		defaultKeyMap.Exit,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, shared.SharedKeyMap.Refresh):
			cmds = append(cmds, m.handleRefreshRequested())
		case key.Matches(msg, defaultKeyMap.Exit):
			cmds = append(cmds, m.handleExit())
		}
	case shared.TickMsg:
		cmds = append(cmds, m.handleRefreshRequested())
	case refreshFeedMsg:
		cmds = append(cmds, m.handleFreshPlays(msg))
	}

	return *m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.hasInitialData {
		return "Loading.."
	}

	frame := frameStyle.Width(m.width)
	t := titleStyle.Width(m.width)
	playCount := len(m.plays)

	if playCount == 0 {
		scoreLine := t.Render(fmt.Sprintf("%s - %s", m.game.Away.Abbreviation, m.game.Home.Abbreviation))
		infoLine := t.Render(fmt.Sprintf("%s", m.game.StartTime.Local().Format("03:04pm")))
		title := titleContainerStyle.Render(lipgloss.JoinVertical(lipgloss.Center, scoreLine, infoLine))

		return frame.Render(title)
	}

	var out strings.Builder

	recentPlay := m.plays[playCount-1]
	scoreLine := t.Render(fmt.Sprintf("%s %d - %d %s", m.game.Away.Abbreviation, recentPlay.AwayScore, recentPlay.HomeScore, m.game.Home.Abbreviation))
	var infoLine string

	if recentPlay.Description == "END GAME" {
		infoLine = t.Render("FINAL")
	} else {
		infoLine = t.Render(fmt.Sprintf("Q%d | %s", recentPlay.Quarter, recentPlay.EndDownDistance))
	}

	title := titleContainerStyle.Render(lipgloss.JoinVertical(lipgloss.Center, scoreLine, infoLine))
	totalHeight := lipgloss.Height(title)
	out.WriteString(title)

	for _, play := range slices.Backward(m.plays) {
		stamp := stampStyle.Render(fmt.Sprintf("%05s", play.Clock))
		availableWidth := m.width - lipgloss.Width(stamp)
		description := playDescriptionStyle.Width(availableWidth).Render(fmt.Sprintf("%s", play.Description))
		line := lipgloss.JoinHorizontal(lipgloss.Top, stamp, description)
		lineHeight := lipgloss.Height(line)
		totalHeight += lineHeight

		if totalHeight > m.height {
			break
		}

		out.WriteString(line)
	}

	return frame.Render(out.String())
}

func (m *Model) handleFreshPlays(next refreshFeedMsg) tea.Cmd {
	m.loading = false

	if next.err != nil {
		return shared.Tick(refreshInterval)
	}

	m.plays = next.plays
	m.hasInitialData = true

	return shared.Tick(refreshInterval)
}

func (m *Model) handleExit() tea.Cmd {
	m.game = api.Scoreboard{}
	m.hasInitialData = false

	return func() tea.Msg {
		return ExitWatchMsg{}
	}
}

func (m *Model) handleRefreshRequested() tea.Cmd {
	if !m.loading {
		m.loading = true

		return func() tea.Msg {
			plays, err := api.GetPlays(m.game.Id)
			return refreshFeedMsg{plays, err}
		}
	}

	return nil
}
