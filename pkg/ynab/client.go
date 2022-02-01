package ynab

import (
	"context"
	"encoding/json"
	"net/http"
)

type APIError struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func (e APIError) Error() string {
	return e.Name
}

type Client struct {
	httpClient *http.Client
	apiToken   string
	BaseURL    string
}

func NewClient(apiToken string) *Client {
	return &Client{
		httpClient: &http.Client{},
		apiToken:   apiToken,
		BaseURL:    "https://api.youneedabudget.com/v1",
	}
}

func (c *Client) Test(ctx context.Context) error {
	req, err := c.newRequest(ctx, "GET", "/user")
	if err != nil {
		return err
	}

	var payload map[string]interface{}

	return c.do(req, &payload)
}

func (c *Client) newRequest(ctx context.Context, method, path string) (*http.Request, error) {
	u := c.BaseURL + path

	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var payload struct {
			Error APIError `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			return err
		}
		return payload.Error
	}

	err = json.NewDecoder(resp.Body).Decode(v)

	return err
}
