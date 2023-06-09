package sntrn

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

const LOGIN_PATH = "/session/userlogin"
const LOGOUT_PATH = "/session/logout"

var logger *log.Logger

// If this changes users can pass a new url as an option to our initializer
var BaseUrl = url.URL{
	Scheme: "https",
	Host:   "senturion.to",
}

type Client struct {
	token string

	host   *url.URL
	client *http.Client
}

type Options struct {
	Host   *url.URL
	Client *http.Client

	LogLevel log.Level
}

func New(opts *Options) *Client {
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Level:           opts.LogLevel,
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})

	if opts.Host == nil {
		opts.Host = &BaseUrl
	}

	if opts.Client == nil {
		opts.Client = defaultHttpClient()
	}

	client := &Client{
		host:   opts.Host,
		client: opts.Client,
	}

	client.index()

	return client
}

// Unlikely that anyone will use this but its here if you already have a token.
func NewWithToken(opts *Options, token string) *Client {
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Level:           opts.LogLevel,
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})

	if opts.Host == nil {
		opts.Host = &BaseUrl
	}

	if opts.Client == nil {
		opts.Client = defaultHttpClient()
	}

	return &Client{
		host:   opts.Host,
		client: opts.Client,
		token:  token,
	}
}

func (c *Client) Close() error {
	if c.token == "" {
		return ErrNotLoggedIn
	}

	u := c.host.ResolveReference(&url.URL{Path: LOGOUT_PATH})

	req, err := c.newRequest("GET", u, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(context.Background(), req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: logout status code not ok", ErrFailedLogout)
	}

	u = c.host.ResolveReference(&url.URL{Path: LOGIN_PATH})

	if resp.Request.URL.String() != u.String() {
		return fmt.Errorf("%w: client not redirected to signin", ErrFailedLogout)
	}

	return nil
}

func (c *Client) newRequest(method string, path *url.URL, body io.Reader) (*http.Request, error) {
	u := c.host.ResolveReference(path)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	req.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36`)
	req.Header.Add("Referer", c.host.String())
	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return nil, err
		}
	}

	defer resp.Body.Close()

	return resp, err
}

func (c *Client) index() error {
	req, err := c.newRequest("GET", c.host, nil)
	if err != nil {
		return err
	}

	_, err = c.do(context.Background(), req, nil)
	if err != nil {
		return err
	}

	// Get the phpsessid and set to token for later
	for _, cookie := range c.client.Jar.Cookies(c.host) {
		if cookie.Name == "PHPSESSID" {
			c.token = cookie.Value
		}
	}

	if c.token == "" {
		return nil
	}

	return nil
}

// End of Client method -----------------------------------------------------------------------------------------------
func defaultHttpClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		DisableCompression: false,
		IdleConnTimeout:    30 * time.Second,
	}

	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		logger.Fatal(err)
	}

	return &http.Client{
		Jar:       jar,
		Transport: tr,
		Timeout:   10 * time.Second,
	}
}
