package ynab

import (
	"context"
	"fmt"
)

type Account struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	Type                string `json:"type,omitempty"`
	OnBudget            bool   `json:"on_budget,omitempty"`
	Closed              bool   `json:"closed,omitempty"`
	Note                string `json:"note,omitempty"`
	Balance             int64  `json:"balance,omitempty"`
	ClearedBalance      int64  `json:"cleared_balance,omitempty"`
	UnclearedBalance    int64  `json:"uncleared_balance,omitempty"`
	TransferPayeeID     string `json:"transfer_payee_id,omitempty"`
	DirectImportLinked  bool   `json:"direct_import_linked,omitempty"`
	DirectImportInError bool   `json:"direct_import_in_error,omitempty"`
	Deleted             bool   `json:"deleted,omitempty"`
}

func (c *Client) Accounts(ctx context.Context, budgetID string) ([]Account, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/budgets/%s/accounts", budgetID))
	if err != nil {
		return nil, err
	}

	var payload struct {
		Data struct {
			Accounts []Account `json:"accounts"`
		} `json:"data"`
	}

	if err := c.do(req, &payload); err != nil {
		return nil, err
	}

	return payload.Data.Accounts, nil
}
