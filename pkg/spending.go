package main

import (
	"context"
	"sort"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/marcusolsson/grafana-ynab-datasource/pkg/ynab"
)

func (d *YNABDataSource) querySpending(ctx context.Context, dsQuery DataSourceQuery, from, to time.Time) backend.DataResponse {
	txs, err := d.client.Transactions(ctx, dsQuery.BudgetID, "", from.Format("2006-01-02"))
	if err != nil {
		return backend.DataResponse{Error: err}
	}

	// categoryGroups, err := d.client.Categories(ctx, dsQuery.BudgetID)
	// if err != nil {
	// 	return backend.DataResponse{Error: err}
	// }

	// categories, err := categoryMappings(categoryGroups)
	// if err != nil {
	// 	return backend.DataResponse{Error: err}
	// }

	// Since the YNAB API returns all transactions from a date until now, we
	// need to filter out any transactions beyond the end of the time range.
	txs, err = filterTransactionsByDate(txs, from, to)
	if err != nil {
		return backend.DataResponse{Error: err}
	}

	txs = filterTransactionsByType(txs, dsQuery.SpendingFilter)

	txs = filterTransactionsByAccountID(txs, dsQuery.AccountIDs)

	period := PeriodDaily

	if dsQuery.Period == "month" {
		period = PeriodMonthly
	} else if dsQuery.Period == "week" {
		period = PeriodWeekly
	}

	groupedTxs := groupTransactions(txs, func(tx ynab.Transaction) string {
		switch dsQuery.AlignBy {
		case "account":
			return tx.AccountID
		case "payee":
			return tx.PayeeID
		// case "category_group":
		// 	return tx.CategoryName
		case "category":
			return tx.CategoryID
		default:
			return ""
		}
	})

	var acc []Measurement
	for _, txs := range groupedTxs {
		measurements, err := Regularize(TimeSeriesTransactions(txs), period, alignTotal, "empty")
		if err != nil {
			return backend.DataResponse{Error: err}
		}
		acc = append(acc, measurements...)
	}

	sort.SliceStable(acc, func(i, j int) bool {
		return acc[i].Time.Before(acc[j].Time)
	})

	frame, err := measurementsToFrame(acc, &data.FillMissing{
		Mode: data.FillModeNull,
	}, dsQuery.AlignBy+"_id", dsQuery.AlignBy+"_name")
	if err != nil {
		return backend.DataResponse{Error: err}
	}

	return backend.DataResponse{
		Frames: data.Frames{frame},
	}
}
