package sntrn

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
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

	requests := []*http.Request{}

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
				continue
			}

			requests = append(requests, req)
		}
	}

	results := []LinksResponse{}
	ch := make(chan LinksResponse)

	var wg sync.WaitGroup

	for _, req := range requests {
		wg.Add(1)
		go func(req *http.Request) {
			defer wg.Done()

			var dr LinksResponse
			_, err = c.do(ctx, req, &dr)
			if err != nil {
				log.Fatal(err)
			}
			ch <- dr
		}(req)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for res := range ch {
		results = append(results, res)
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
			break // Occassionally this is because the servers have a season 0
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
