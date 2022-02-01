package main

import (
	"context"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/marcusolsson/grafana-ynab-datasource/pkg/ynab"
)

func (d *YNABDataSource) queryTransactions(ctx context.Context, dsQuery DataSourceQuery, from, to time.Time) backend.DataResponse {
	txs, err := d.client.Transactions(ctx, dsQuery.BudgetID, "", from.Format("2006-01-02"))
	if err != nil {
		return backend.DataResponse{Error: err}
	}

	txs = filterTransactionsByAccountID(txs, dsQuery.AccountIDs)

	categoryGroups, err := d.client.Categories(ctx, dsQuery.BudgetID)
	if err != nil {
		return backend.DataResponse{Error: err}
	}

	frame, err := FrameFromTransactions(txs, categoryGroups, dsQuery, from, to)
	if err != nil {
		return backend.DataResponse{Error: err}
	}

	frame.Meta = &data.FrameMeta{
		PreferredVisualization: data.VisTypeTable,
	}

	return backend.DataResponse{
		Frames: data.Frames{frame},
	}
}

func FrameFromTransactions(transactions []ynab.Transaction, categoryGroups []ynab.CategoryGroup, dsQuery DataSourceQuery, from, to time.Time) (*data.Frame, error) {
	categories, err := categoryMappings(categoryGroups)
	if err != nil {
		return nil, err
	}

	// Since the YNAB API returns all transactions from a date until now, we
	// need to filter out any transactions beyond the end of the time range.
	txs, err := filterTransactionsByDate(transactions, from, to)
	if err != nil {
		return nil, err
	}

	txs = filterTransactionsByType(txs, dsQuery.TransactionFilter)

	frame, err := transactionsFrame(txs, categories)
	if err != nil {
		return nil, err
	}

	return frame, nil
}
