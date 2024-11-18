package shared

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	TickMsg time.Time

	keyMap struct {
		Refresh key.Binding
	}
)

var (
	SharedKeyMap = keyMap{
		Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
	}
)

func Tick(duration time.Duration) tea.Cmd {
	return tea.Every(duration, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
