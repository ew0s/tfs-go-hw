package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"trade-bot/configs"
)

type ClientActions interface {
	NewRequest(method, path, jwtToken string, body interface{}) (*http.Request, error)
	Do(req *http.Request, typ interface{}) (*http.Response, error)
}

type Client struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *http.Client
}

func NewClient(config configs.ClientConfiguration) (*Client, error) {
	urlValue, err := url.Parse(config.URL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		BaseURL:    urlValue,
		UserAgent:  "go client user agent",
		httpClient: http.DefaultClient,
	}

	return c, nil
}

func (c *Client) NewRequest(method, path, jwtToken string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	if jwtToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, typ interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(typ)
	return resp, err
}
