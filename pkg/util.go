package main

import (
	"regexp"
	"strconv"
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

func transactionsFrame(transactions []ynab.Transaction, groups map[string]ynab.CategoryGroup, unit string) (*data.Frame, error) {

	amountField := data.NewField("amount", nil, []float64{})
	amountField.Config = &data.FieldConfig{
		Unit: unit,
	}

	frame := data.NewFrame(
		"transactions",
		data.NewField("time", nil, []time.Time{}),
		data.NewField("account", nil, []string{}),
		data.NewField("payee", nil, []string{}),
		amountField,
		data.NewField("memo", nil, []string{}),
		data.NewField("category", nil, []string{}),
		data.NewField("category_group", nil, []*string{}),
		data.NewField("stock", nil, []string{}),
		data.NewField("quantity", nil, []int{}),
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

		stock, quantity, err := findStockPositions(tx.Memo)

		if err == nil && stock != "" && quantity > 0 {
			frame.AppendRow(
				date,
				tx.AccountName,
				tx.PayeeName,
				amount,
				tx.Memo,
				tx.CategoryName,
				categoryGroup,
				stock,
				quantity,
			)
		} else {
			frame.AppendRow(
				date,
				tx.AccountName,
				tx.PayeeName,
				amount,
				tx.Memo,
				tx.CategoryName,
				categoryGroup,
				"",
				0,
			)
		}
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

func measurementsToFrame(measurements []Measurement, fillMissing *data.FillMissing, idLabel, nameLabel string, unit string) (*data.Frame, error) {
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
			Unit:              unit,
		}
	}

	return wideFrame, nil
}

func convertCurrencyCode(isoCode string) string {
	switch isoCode {
	case "USD", "GBP", "EUR", "JPY", "RUB", "UAH", "BRL", "DKK", "ISK", "NOK", "SEK", "CZK", "CHF", "PLN", "ZAR", "INR", "KRW", "IDR", "PHP", "VND":
		return "currency" + isoCode
	}
	return "locale"
}

func findStockPositions(memo string) (string, int, error) {

	//Matches a label=value structure, ie "stock=AAPL, quantity=200"
	regexStockPosition := regexp.MustCompile(`(?P<key>\b[A-Za-z.]+)=(?P<value>\w+)`)

	var stockKey = "stock"
	var quantityKey = "quantity"
	var stock string
	var quantity int

	for _, match := range regexStockPosition.FindAllStringSubmatch(memo, -1) {
		if match[1] == stockKey {
			stock = match[2]
		}
		if match[1] == quantityKey {
			quantityValue, err := strconv.Atoi(match[2])
			if err != nil {
				return "", 0, err
			}
			quantity = quantityValue
		}
	}

	return stock, quantity, nil
}
