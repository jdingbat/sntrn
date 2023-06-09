package sntrn

import (
	"context"
	"strings"
)

type LinksResponse struct {
	Jwplayer []struct {
		File  string `json:"file"`
		Type  string `json:"type"`
		Label string `json:"label"`
	} `json:"jwplayer"`

	Dl   string `json:"dl"`
	DlHd string `json:"dl_hd"`

	Server struct {
		List     map[string]string `json:"list"`
		Host     string            `json:"host"`
		Selected int               `json:"selected"`
		Scheme   string            `json:"scheme"`
		P720     string            `json:"720"`
		P1080    string            `json:"1080"`
	} `json:"server"`
}

func (c *Client) Links(ctx context.Context, sr SearchResponse) ([]LinksResponse, error) {
	if strings.HasPrefix(sr.Link, "/movies") {
		return c.getMovieLinks(ctx, sr)
	}

	return c.getSeriesLinks(ctx, sr)
}
