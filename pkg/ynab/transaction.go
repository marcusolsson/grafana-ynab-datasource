package ynab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type Transaction struct {
	ID                    string        `json:"id,omitempty"`
	Date                  string        `json:"date,omitempty"`
	Amount                int64         `json:"amount,omitempty"`
	Memo                  string        `json:"memo,omitempty"`
	Cleared               string        `json:"cleared,omitempty"`
	Approved              bool          `json:"approved,omitempty"`
	FlagColor             string        `json:"flag_color,omitempty"`
	AccountID             string        `json:"account_id,omitempty"`
	PayeeID               string        `json:"payee_id,omitempty"`
	CategoryID            string        `json:"category_id,omitempty"`
	TransferAccountID     string        `json:"transfer_account_id,omitempty"`
	TransferTransactionID string        `json:"transfer_transaction_id,omitempty"`
	MatchedTransactionID  string        `json:"matched_transaction_id,omitempty"`
	ImportID              string        `json:"import_id,omitempty"`
	Deleted               bool          `json:"deleted,omitempty"`
	AccountName           string        `json:"account_name,omitempty"`
	PayeeName             string        `json:"payee_name,omitempty"`
	CategoryName          string        `json:"category_name,omitempty"`
	Subtransactions       []Transaction `json:"subtransactions,omitempty"`
}

type Subtransaction struct {
	ID                    string `json:"id,omitempty"`
	TransactionID         string `json:"transaction_id,omitempty"`
	Amount                int64  `json:"amount,omitempty"`
	Memo                  string `json:"memo,omitempty"`
	PayeeID               string `json:"payee_id,omitempty"`
	PayeeName             string `json:"payee_name,omitempty"`
	CategoryID            string `json:"category_id,omitempty"`
	CategoryName          string `json:"category_name,omitempty"`
	TransferAccountID     string `json:"transfer_account_id,omitempty"`
	TransferTransactionID string `json:"transfer_transaction_id,omitempty"`
	Deleted               bool   `json:"deleted,omitempty"`
}

func (c *Client) Transactions(ctx context.Context, budgetID string, sinceDate string) ([]Transaction, error) {
	backend.Logger.Error("YNABClient.TransactionsForAccount()", "budgetID", budgetID, "sinceDate", sinceDate)

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+fmt.Sprintf("/budgets/%s/transactions", budgetID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	q := req.URL.Query()
	q.Set("since_date", sinceDate)
	req.URL.RawQuery = q.Encode()

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
			Transactions []Transaction `json:"transactions"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Data.Transactions, nil
}

func (c *Client) TransactionsForAccount(ctx context.Context, budgetID, accountID string, sinceDate string) ([]Transaction, error) {
	backend.Logger.Error("YNABClient.TransactionsForAccount()", "budgetID", budgetID, "accountID", accountID, "sinceDate", sinceDate)

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+fmt.Sprintf("/budgets/%s/accounts/%s/transactions", budgetID, accountID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	q := req.URL.Query()
	q.Set("since_date", sinceDate)
	req.URL.RawQuery = q.Encode()

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
			Transactions []Transaction `json:"transactions"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Data.Transactions, nil
}
