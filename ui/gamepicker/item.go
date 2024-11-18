package gamepicker

import (
	"fmt"

	"github.com/s992/nfl/api"
)

type (
	item struct {
		api.Scoreboard
	}
)

func (i item) Title() string       { return i.Name }
func (i item) FilterValue() string { return i.Name }
func (i item) Description() string {
	if i.State == api.GameStatePre {
		return fmt.Sprintf("Kickoff @ %s", i.StartTime.Local().Format("03:04pm"))
	}

	time := ""
	if i.State == api.GameStatePost {
		time = "FINAL   "
	} else {
		time = fmt.Sprintf("Q%d %05s", i.Quarter, i.Clock)
	}

	return fmt.Sprintf("%s | %3s %2d - %-2d %-3s", time, i.Away.Abbreviation, i.Away.Score, i.Home.Score, i.Home.Abbreviation)
}
