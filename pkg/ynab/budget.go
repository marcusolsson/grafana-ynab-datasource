package ynab

import (
	"context"
)

type Budget struct {
	ID       string    `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Accounts []Account `json:"accounts,omitempty"`
}

func (c *Client) Budgets(ctx context.Context, includeAccounts bool) ([]Budget, error) {
	req, err := c.newRequest(ctx, "GET", "/budgets")
	if err != nil {
		return nil, err
	}

	if includeAccounts {
		req.URL.RawQuery = "include_accounts=true"
	}

	var payload struct {
		Data struct {
			Budgets []Budget `json:"budgets"`
		} `json:"data"`
	}

	if err := c.do(req, &payload); err != nil {
		return nil, err
	}

	return payload.Data.Budgets, nil
}
