package sntrn

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"time"
)

type SearchResponse struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Year       string `json:"year"`
	Genre      string `json:"genre"`
	ImdbRating string `json:"imdb_rating"`
	Poster     string `json:"poster"`
	Link       string `json:"link"`
}

func (c *Client) Search(ctx context.Context, query string) ([]SearchResponse, error) {
	if len(query) < 2 {
		return nil, errors.New("test")
	}

	params := url.Values{
		"q": {query},
		"_": {strconv.FormatInt(time.Now().UnixMilli(), 10)},
	}

	u := c.host.ResolveReference(&url.URL{Path: "search/auto"})
	u.RawQuery = params.Encode()

	req, err := c.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var sr []SearchResponse
	_, err = c.do(ctx, req, &sr)
	if err != nil {
		return nil, err
	}

	return sr, nil
}
