package api

import (
	"io"
	"net/http"
	"time"
)

func fetch(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func parseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04Z", date)
}
