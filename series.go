package sntrn

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type EpisodeData struct {
	EpisodeNumber string `json:"episode_number"`
	Title         string `json:"title"`
	Plot          string `json:"plot"`
	Poster        string `json:"poster"`
	ReleaseDate   string `json:"release_date"`
}

func (c *Client) getSeriesLinks(ctx context.Context, sr SearchResponse) ([]LinksResponse, error) {
	seasons, err := c.getSeasons(ctx, sr)
	if err != nil {
		return nil, err
	}

	results := []LinksResponse{}

	for season, episodes := range seasons {
		for _, episode := range episodes {
			u := c.host.ResolveReference(&url.URL{Path: "/series/getTvLink"})
			u.RawQuery = url.Values{
				"id":    {sr.Id},
				"token": {c.token},
				"s":     {strconv.Itoa(season)},
				"e":     {episode},
				"oPid":  {},
				"_":     {strconv.FormatInt(time.Now().UnixMilli(), 10)},
			}.Encode()

			req, err := c.newRequest("GET", u, nil)
			if err != nil {
				return nil, err
			}

			var dr LinksResponse
			_, err = c.do(ctx, req, &dr)
			if err != nil {
				return nil, err
			}

			results = append(results, dr)
		}
	}

	return results, nil
}

func (c *Client) getSeasons(ctx context.Context, sr SearchResponse) (map[int][]string, error) {
	episodes := map[int][]string{}

	season := 1

	for {
		u := c.host.ResolveReference(&url.URL{Path: "/series/season"})
		u.RawQuery = url.Values{
			"id":    {sr.Id},
			"s":     {strconv.Itoa(season)},
			"token": {c.token},
			"_":     {strconv.FormatInt(time.Now().UnixMilli(), 10)},
		}.Encode()

		req, err := c.newRequest("GET", u, nil)
		if err != nil {
			return nil, err
		}

		var ep []EpisodeData
		resp, err := c.do(ctx, req, &ep)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusInternalServerError {
			break
		}

		episodes[season] = []string{}
		for _, episode := range ep {
			episodes[season] = append(episodes[season], episode.EpisodeNumber)
		}

		season += 1
	}

	return episodes, nil
}
