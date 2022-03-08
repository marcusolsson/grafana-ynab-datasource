package ynab

import (
	"context"
	"fmt"
)

type Budget struct {
	ID             string         `json:"id,omitempty"`
	Name           string         `json:"name,omitempty"`
	Accounts       []Account      `json:"accounts,omitempty"`
	CurrencyFormat CurrencyFormat `json:"currency_format"`
}

type CurrencyFormat struct {
	ISOCode       string `json:"iso_code"`
	DecimalDigits int    `json:"decimal_digits"`
}

func (c *Client) Budget(ctx context.Context, budgetID string) (Budget, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/budgets/%s", budgetID))
	if err != nil {
		return Budget{}, err
	}

	var payload struct {
		Data struct {
			Budget Budget `json:"budget"`
		} `json:"data"`
	}

	if err := c.do(req, &payload); err != nil {
		return Budget{}, err
	}

	return payload.Data.Budget, nil
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
