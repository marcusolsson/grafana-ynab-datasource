package ynab

import (
	"context"
	"fmt"
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

func (c *Client) Transactions(ctx context.Context, budgetID, accountID, sinceDate string) ([]Transaction, error) {
	if accountID != "" {
		return c.transactionsForAccount(ctx, budgetID, accountID, sinceDate)
	}

	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/budgets/%s/transactions", budgetID))
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()

	if sinceDate != "" {
		q.Set("since_date", sinceDate)
		req.URL.RawQuery = q.Encode()
	}

	var payload struct {
		Data struct {
			Transactions []Transaction `json:"transactions"`
		} `json:"data"`
	}

	if err := c.do(req, &payload); err != nil {
		return nil, err
	}

	return payload.Data.Transactions, nil
}

func (c *Client) transactionsForAccount(ctx context.Context, budgetID, accountID string, sinceDate string) ([]Transaction, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/budgets/%s/accounts/%s/transactions", budgetID, accountID))
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if sinceDate != "" {
		q.Set("since_date", sinceDate)
		req.URL.RawQuery = q.Encode()
	}

	var payload struct {
		Data struct {
			Transactions []Transaction `json:"transactions"`
		} `json:"data"`
	}

	if err := c.do(req, &payload); err != nil {
		return nil, err
	}

	return payload.Data.Transactions, nil
}
