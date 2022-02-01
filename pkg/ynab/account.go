package ynab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
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
	backend.Logger.Error("YNABClient.Accounts()", "budgetID", budgetID)

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+fmt.Sprintf("/budgets/%s/accounts", budgetID), nil)
	if err != nil {
		return nil, err
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
			Accounts []Account `json:"accounts"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Data.Accounts, nil
}
