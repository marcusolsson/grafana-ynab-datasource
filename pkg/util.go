package main

import (
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/marcusolsson/grafana-ynab-datasource/pkg/ynab"
)

// filterTransactionsByDate returns all transactions made between the from and to timestamps.
func filterTransactionsByDate(txs []ynab.Transaction, from, to time.Time) ([]ynab.Transaction, error) {
	var res []ynab.Transaction

	for _, tx := range txs {
		t, err := time.Parse("2006-01-02", tx.Date)
		if err != nil {
			return nil, err
		}

		if !t.Before(from) && !t.After(to) {
			res = append(res, tx)
		}
	}

	return res, nil
}

func filterTransactionsByType(txs []ynab.Transaction, filter string) []ynab.Transaction {
	var res []ynab.Transaction

	for _, tx := range txs {
		switch filter {
		case "spending":
			if tx.Amount < 0 {
				tx.Amount = -tx.Amount
				res = append(res, tx)
			}
		case "income":
			if tx.Amount > 0 {
				res = append(res, tx)
			}
		case "all":
			res = append(res, tx)
		}
	}

	return res
}

func filterTransactionsByAccountID(txs []ynab.Transaction, ids []string) []ynab.Transaction {
	if len(ids) == 0 {
		return txs
	}

	var res []ynab.Transaction

	lookup := map[string]bool{}
	for _, id := range ids {
		lookup[id] = true
	}

	for _, tx := range txs {
		if _, ok := lookup[tx.AccountID]; ok {
			res = append(res, tx)
		}
	}

	return res
}

// categoryMappings returns mapping between a category and its category group.
func categoryMappings(categoryGroups []ynab.CategoryGroup) (map[string]ynab.CategoryGroup, error) {
	res := make(map[string]ynab.CategoryGroup)

	for _, categoryGroup := range categoryGroups {
		for _, category := range categoryGroup.Categories {
			res[category.ID] = categoryGroup
		}
	}

	return res, nil
}

func transactionsFrame(transactions []ynab.Transaction, groups map[string]ynab.CategoryGroup) (*data.Frame, error) {
	frame := data.NewFrame(
		"transactions",
		data.NewField("time", nil, []time.Time{}),
		data.NewField("account", nil, []string{}),
		data.NewField("payee", nil, []string{}),
		data.NewField("amount", nil, []float64{}),
		data.NewField("memo", nil, []string{}),
		data.NewField("category", nil, []string{}),
		data.NewField("category_group", nil, []*string{}),
	)

	for _, tx := range transactions {
		date, err := time.Parse("2006-01-02", tx.Date)
		if err != nil {
			return nil, err
		}

		var categoryGroup *string

		// Not all categories belongs to a category group.
		if group, ok := groups[tx.CategoryID]; ok {
			val := group.Name
			categoryGroup = &val
		}

		// YNAB returns amounts in "milliunits". To be able to format
		// them as proper currencies, we need to convert them into
		// regular "units".
		amount := float64(tx.Amount) / 1000.0

		frame.AppendRow(
			date,
			tx.AccountName,
			tx.PayeeName,
			amount,
			tx.Memo,
			tx.CategoryName,
			categoryGroup,
		)
	}

	return frame, nil
}

func groupTransactionsByAccountID(txs []ynab.Transaction) map[string][]ynab.Transaction {
	res := make(map[string][]ynab.Transaction)

	for _, tx := range txs {
		if _, ok := res[tx.AccountID]; !ok {
			res[tx.AccountID] = []ynab.Transaction{}
		}
		res[tx.AccountID] = append(res[tx.AccountID], tx)
	}

	return res
}

func groupTransactions(txs []ynab.Transaction, fn func(tx ynab.Transaction) string) map[string][]ynab.Transaction {
	res := make(map[string][]ynab.Transaction)

	for _, tx := range txs {
		key := fn(tx)

		if _, ok := res[key]; !ok {
			res[key] = []ynab.Transaction{}
		}

		res[key] = append(res[key], tx)
	}

	return res
}

func netWorthFromFrame(frame *data.Frame) (*data.Frame, error) {
	networth := data.NewFrame(
		"net_worth",
		data.NewField("time", nil, []time.Time{}),
		data.NewField("debts", nil, []*float64{}),
		data.NewField("assets", nil, []*float64{}),
		data.NewField("net_worth", nil, []*float64{}),
	)

	for i := 0; i < frame.Rows(); i++ {
		timeField := frame.Fields[0]
		valueFields := frame.Fields[1:]

		var income *float64
		var spending *float64
		net := new(float64)

		for _, field := range valueFields {
			val := field.CopyAt(i).(*float64)
			if val == nil {
				continue
			}

			if *val > 0 {
				if income == nil {
					income = val
				} else {
					*income += *val
				}

			}
			if *val < 0 {
				if spending == nil {
					spending = val
				} else {
					*spending += *val
				}

			}
		}

		if income != nil {
			if spending != nil {
				*spending = -*spending
				*net = *income - *spending
			} else {
				*net = *income
			}
		}

		networth.AppendRow(timeField.At(i), spending, income, net)
	}

	return networth, nil
}

func measurementsToFrame(measurements []Measurement, fillMissing *data.FillMissing, idLabel, nameLabel string) (*data.Frame, error) {
	longFrame := data.NewFrame(
		"balance",
		data.NewField("time", nil, []time.Time{}),
		data.NewField("amount", nil, []float64{}),
		data.NewField("id", nil, []string{}),
		data.NewField("name", nil, []string{}),
	)

	for _, measurement := range measurements {
		longFrame.AppendRow(measurement.Time, float64(measurement.Value)/1000.0, measurement.Labels[idLabel], measurement.Labels[nameLabel])
	}

	wideFrame, err := data.LongToWide(longFrame, fillMissing)
	if err != nil {
		return nil, err
	}

	for _, field := range wideFrame.Fields {
		displayName := field.Labels["name"]
		field.Config = &data.FieldConfig{
			DisplayName:       displayName,
			DisplayNameFromDS: displayName,
		}
	}

	return wideFrame, nil
}
