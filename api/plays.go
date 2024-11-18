package api

import (
	"encoding/json"
	"fmt"
)

type (
	Play struct {
		Description       string
		AwayScore         uint
		HomeScore         uint
		Quarter           uint
		Clock             string
		StartDownDistance string
		EndDownDistance   string
	}
)

type (
	feed struct {
		Count     uint
		PageIndex uint
		PageSize  uint
		PageCount uint
		Items     []feedEntry
	}

	feedEntry struct {
		Text      string
		AwayScore uint
		HomeScore uint
		Period    struct {
			Number uint
		}
		Clock struct {
			DisplayValue string
		}
		Start playInfo
		End   playInfo
	}

	playInfo struct {
		DownDistanceText string
	}
)

const (
	pageSize = 500
)

func GetPlays(gameId string) ([]Play, error) {
	var entries []feedEntry
	firstPage, err := fetchFeedPage(gameId, 1)
	if err != nil {
		return nil, err
	}

	entries = append(entries, firstPage.Items...)

	page := firstPage.PageIndex + 1
	for page < firstPage.PageCount {
		nextPage, err := fetchFeedPage(gameId, page)
		if err != nil {
			return nil, err
		}

		entries = append(entries, nextPage.Items...)
		page++
	}

	var plays []Play
	for _, entry := range entries {
		plays = append(plays, Play{
			Description:       entry.Text,
			AwayScore:         entry.AwayScore,
			HomeScore:         entry.HomeScore,
			Quarter:           entry.Period.Number,
			Clock:             entry.Clock.DisplayValue,
			StartDownDistance: entry.Start.DownDistanceText,
			EndDownDistance:   entry.End.DownDistanceText,
		})
	}

	return plays, nil
}

func fetchFeedPage(gameId string, page uint) (*feed, error) {
	url := fmt.Sprintf("https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/events/%s/competitions/%s/plays/?limit=%d&page=%d", gameId, gameId, pageSize, page)
	res, err := fetch(url)
	if err != nil {
		return nil, err
	}

	var resFeed *feed
	err = json.Unmarshal(res, &resFeed)
	if err != nil {
		return nil, err
	}

	return resFeed, nil
}
