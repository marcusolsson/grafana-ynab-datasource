package ynab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type Budget struct {
	ID       string    `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Accounts []Account `json:"accounts,omitempty"`
}

func (c *Client) Budgets(ctx context.Context, includeAccounts bool) ([]Budget, error) {
	backend.Logger.Error("YNABClient.Budgets()", "includeAccounts", includeAccounts)

	// return []Budget{
	// 	{ID: "foo", Name: "Foo", Accounts: []Account{{ID: "foo-checking", Name: "Checking"}, {ID: "foo-buffer", Name: "Buffer"}}},
	// 	{ID: "bar", Name: "Bar", Accounts: []Account{{ID: "bar-ica", Name: "ICA Banken"}}},
	// }, nil

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/budgets", nil)
	if err != nil {
		return nil, err
	}

	if includeAccounts {
		req.URL.RawQuery = "include_accounts=true"
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var payload struct {
		Data struct {
			Budgets []Budget `json:"budgets"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Data.Budgets, nil
}
