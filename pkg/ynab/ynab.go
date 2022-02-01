package ynab

import (
	"context"
	"errors"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type Client struct {
	httpClient *http.Client
	apiToken   string
	baseURL    string
}

func NewClient(apiToken string) *Client {
	return &Client{
		httpClient: &http.Client{},
		apiToken:   apiToken,
		baseURL:    "https://api.youneedabudget.com/v1",
	}
}

func (c *Client) Test(ctx context.Context) (int, error) {
	backend.Logger.Error("YNABClient.Test()")

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/user", nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return resp.StatusCode, errors.New("unexpected status code")
	}

	return 200, nil
}
