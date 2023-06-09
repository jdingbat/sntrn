package sntrn

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

func (c *Client) getMovieLinks(ctx context.Context, sr SearchResponse) ([]LinksResponse, error) {
	u := c.host.ResolveReference(&url.URL{Path: "/movies/getMovieLink"})
	u.RawQuery = url.Values{
		"id":    {sr.Id},
		"token": {c.token},
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

	return []LinksResponse{dr}, nil
}
