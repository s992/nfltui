package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type (
	Scoreboard struct {
		Id        string
		StartTime time.Time
		Name      string
		ShortName string
		Clock     string
		Quarter   uint
		Drive     string
		LastPlay  string
		State     GameState
		Home      Competitor
		Away      Competitor
	}

	Competitor struct {
		DisplayName      string
		ShortDisplayName string
		Abbreviation     string
		Score            uint
	}

	GameState string
)

const (
	GameStatePre        GameState = "pre"
	GameStatePost       GameState = "post"
	GameStateInProgress GameState = "in"
	GameStateUnknown    GameState = "unknown"
)

type (
	scoreboardResponse struct {
		Events []scoreboardEvent
	}

	scoreboardEvent struct {
		Id           string
		Date         string
		Name         string
		ShortName    string
		Competitions []competition
	}

	competition struct {
		Competitors []competitor
		Situation   struct {
			LastPlay struct {
				Text string
			}
			Probability struct {
				TiePercentage     float32
				HomeWinPercentage float32
				AwayWinPercentage float32
			}
			Drive struct {
				Description string
			}
		}
		Status struct {
			DisplayClock string
			Period       uint
			Type         struct {
				State string
			}
		}
	}

	competitor struct {
		HomeAway string
		Team     struct {
			DisplayName      string
			ShortDisplayName string
			Abbreviation     string
		}
		Score      string
		Linescores []struct {
			Value float32
		}
	}
)

func GetScoreboards(date time.Time) ([]Scoreboard, error) {
	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard?dates=%s", date.Format("20060102"))
	res, err := fetch(url)
	if err != nil {
		return nil, err
	}

	var board scoreboardResponse
	err = json.Unmarshal(res, &board)
	if err != nil {
		return nil, err
	}

	var boards []Scoreboard

	for _, event := range board.Events {
		comp := event.Competitions[0]
		home, away := getHomeAway(comp.Competitors)
		startTime, _ := parseDate(event.Date)
		homeScore, _ := strconv.Atoi(home.Score)
		awayScore, _ := strconv.Atoi(away.Score)
		state := GameStateUnknown

		switch comp.Status.Type.State {
		case "pre":
			state = GameStatePre
		case "post":
			state = GameStatePost
		case "in":
			state = GameStateInProgress
		}

		boards = append(boards, Scoreboard{
			Id:        event.Id,
			StartTime: startTime,
			Name:      event.Name,
			ShortName: event.ShortName,
			Clock:     comp.Status.DisplayClock,
			Quarter:   comp.Status.Period,
			Drive:     comp.Situation.Drive.Description,
			LastPlay:  comp.Situation.LastPlay.Text,
			State:     state,
			Home: Competitor{
				DisplayName:      home.Team.DisplayName,
				ShortDisplayName: home.Team.ShortDisplayName,
				Abbreviation:     home.Team.Abbreviation,
				Score:            uint(homeScore),
			},
			Away: Competitor{
				DisplayName:      away.Team.DisplayName,
				ShortDisplayName: away.Team.ShortDisplayName,
				Abbreviation:     away.Team.Abbreviation,
				Score:            uint(awayScore),
			},
		})
	}

	return boards, nil
}

func getHomeAway(c []competitor) (competitor, competitor) {
	if c[0].HomeAway == "home" {
		return c[0], c[1]
	}

	return c[1], c[0]
}
