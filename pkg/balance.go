package main

import (
	"context"
	"sort"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/marcusolsson/grafana-ynab-datasource/pkg/ynab"
)

func (d *YNABDataSource) queryBalance(ctx context.Context, dsQuery DataSourceQuery, from, to time.Time) backend.DataResponse {
	period := PeriodDaily

	if dsQuery.Period == "month" {
		period = PeriodMonthly
	} else if dsQuery.Period == "week" {
		period = PeriodWeekly
	}

	txs, err := d.client.Transactions(ctx, dsQuery.BudgetID, "", "")
	if err != nil {
		return backend.DataResponse{Error: err}
	}

	txs = filterTransactionsByAccountID(txs, dsQuery.AccountIDs)

	groupedTxs := groupTransactionsByAccountID(txs)

	var acc []Measurement
	for _, txs := range groupedTxs {
		balances, err := ynab.RunningBalance(txs)
		if err != nil {
			return backend.DataResponse{Error: err}
		}

		measurements, err := Regularize(TimeSeriesBalance(balances), period, alignLast, "last")
		if err != nil {
			return backend.DataResponse{Error: err}
		}
		acc = append(acc, measurements...)
	}

	sort.SliceStable(acc, func(i, j int) bool {
		return acc[i].Time.Before(acc[j].Time)
	})

	frame, err := measurementsToFrame(acc, &data.FillMissing{Mode: data.FillModePrevious}, "account_id", "account_name")
	if err != nil {
		return backend.DataResponse{Error: err}
	}

	return backend.DataResponse{
		Frames: data.Frames{frame},
	}
}
