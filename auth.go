package sntrn

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) Login(ctx context.Context, username, password string) error {
	if len(password) < 8 {
		return fmt.Errorf("%w: password doesn't match requirements", ErrFailedLogin)
	}

	data := url.Values{
		"username": {username},
		"password": {password},
		"remember": {"0"},
	}

	u := c.host.ResolveReference(&url.URL{Path: LOGIN_PATH})

	// This basically acts like `PostForm`
	req, err := c.newRequest(http.MethodPost, u, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedLogin, err)
	}

	resp, err := c.do(ctx, req, nil)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedLogin, err)
	}

	if resp.Request.URL.Path == LOGIN_PATH {
		return fmt.Errorf("%w: failed to login, bad username or password", ErrFailedLogin)
	}

	return nil
}
