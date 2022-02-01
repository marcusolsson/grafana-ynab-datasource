package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/marcusolsson/grafana-ynab-datasource/pkg/ynab"
)

type SelectableValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

var (
	_ backend.QueryDataHandler    = (*YNABDatasource)(nil)
	_ backend.CheckHealthHandler  = (*YNABDatasource)(nil)
	_ backend.CallResourceHandler = (*YNABDatasource)(nil)
)

// queryModel represents the query that's constructed by the query editor and
// sent by the dashboard.
type queryModel struct {
	BudgetID  string `json:"budgetId"`
	AccountID string `json:"accountId"`
	AlignBy   string `json:"alignBy"`
	Period    string `json:"period"`
	QueryType string `json:"queryType"`
}

type YNABDatasource struct {
	client *ynab.Client
}

// NewYNABDatasource creates a new datasource instance.
func NewYNABDatasource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	apiToken := settings.DecryptedSecureJSONData["apiToken"]

	return &YNABDatasource{
		client: ynab.NewClient(apiToken),
	}, nil
}

func (d *YNABDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()
	for _, q := range req.Queries {
		response.Responses[q.RefID] = d.query(ctx, req.PluginContext, q)
	}
	return response, nil
}

func (d *YNABDatasource) getTransactions(ctx context.Context, budgetID, accountID, sinceDate string) ([]ynab.Transaction, error) {
	var res []ynab.Transaction
	var err error

	if accountID == "" {
		res, err = d.client.Transactions(ctx, budgetID, sinceDate)
		if err != nil {
			return nil, err
		}
	} else {
		res, err = d.client.TransactionsForAccount(ctx, budgetID, accountID, sinceDate)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (d *YNABDatasource) filterTransactionsByDate(txs []ynab.Transaction, from, to time.Time) ([]ynab.Transaction, error) {
	var filteredTransactions []ynab.Transaction

	for _, tx := range txs {
		t, err := time.Parse("2006-01-02", tx.Date)
		if err != nil {
			backend.Logger.Error("unable to parse date", "date", tx.Date, "err", err.Error())
			continue
		}

		if !t.Before(from) && !t.After(to) {
			filteredTransactions = append(filteredTransactions, tx)
		}
	}

	return filteredTransactions, nil
}

func (d *YNABDatasource) categoryGroupLookup(ctx context.Context, budgetID string) (map[string]ynab.CategoryGroup, error) {
	cats, err := d.client.Categories(ctx, budgetID)
	if err != nil {
		return nil, err
	}

	categoryGroups := make(map[string]ynab.CategoryGroup)
	for _, g := range cats {
		for _, c := range g.Categories {
			categoryGroups[c.ID] = g
		}
	}

	return categoryGroups, nil
}

func (d *YNABDatasource) query(ctx context.Context, pluginCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var response backend.DataResponse

	var qm queryModel
	if err := json.Unmarshal(query.JSON, &qm); err != nil {
		response.Error = err
		return response
	}

	if qm.BudgetID == "" {
		return response
	}

	var err error
	var txs []ynab.Transaction

	txs, err = d.getTransactions(ctx, qm.BudgetID, qm.AccountID, query.TimeRange.From.Format("2006-01-02"))
	if err != nil {
		response.Error = err
		return response
	}

	txs, err = d.filterTransactionsByDate(txs, query.TimeRange.From, query.TimeRange.To)
	if err != nil {
		response.Error = err
		return response
	}

	groups, err := d.categoryGroupLookup(ctx, qm.BudgetID)
	if err != nil {
		response.Error = err
		return response
	}

	frame, err := transactionsToWideFrame(txs, groups, qm.AlignBy, qm.QueryType)
	if err != nil {
		backend.Logger.Error("error", "err", err.Error())
		response.Error = err
		return response
	}

	frame, err = alignByPeriod(frame, qm.Period)
	if err != nil {
		response.Error = err
		return response
	}

	response.Frames = data.Frames{frame}

	return response
}

func alignByPeriod(frame *data.Frame, period string) (*data.Frame, error) {
	switch period {
	case "day":
		return frame, nil
	case "month":
		timeField, timeIdx := frame.FieldByName("time")

		aligned := make(map[string][]int)
		for i := 0; i < frame.Rows(); i++ {
			t := timeField.At(i).(time.Time)

			truncated := t.Format("2006-01")

			aligned[truncated] = append(aligned[truncated], i)
		}

		res := frame.EmptyCopy()

		for month, ids := range aligned {
			t, err := time.Parse("2006-01", month)
			if err != nil {
				continue
			}
			res.Fields[timeIdx].Append(t)

			// Reduce each field.
			for i, field := range res.Fields {
				backend.Logger.Error("reducing field", "type", field.Type())
				if field.Type() == data.FieldTypeNullableFloat64 {
					var total float64
					for _, id := range ids {
						val, _ := frame.Fields[i].ConcreteAt(id)
						if val != nil {
							total += val.(float64)
						}
					}
					res.Fields[i].Append(&total)
				}
			}
		}

		return res, nil
	}

	return nil, errors.New("unsupported period")
}

// CheckHealth runs when the user presses "Save & Test" in the data source settings. Tests that the user can access
// their workspace.
func (d *YNABDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	status, err := d.client.Test(ctx)
	if err != nil {
		switch status {
		case http.StatusUnauthorized:
			return &backend.CheckHealthResult{
				Status:  backend.HealthStatusError,
				Message: "Couldn't authorize",
			}, nil
		default:
			return &backend.CheckHealthResult{
				Status:  backend.HealthStatusError,
				Message: fmt.Sprintf("Received an unexpected status code: %d", status),
			}, nil
		}
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Success",
	}, nil
}

// CallResources exposes a REST API with support operations for the query editor.
func (d *YNABDatasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	path := strings.Split(req.Path, "/")

	if len(path) == 0 {
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusNotFound,
		})
	}

	if len(path) == 1 && path[0] == "budgets" {
		budgets, err := d.client.Budgets(ctx, true)
		if err != nil {
			return err
		}

		b, err := json.Marshal(budgets)
		if err != nil {
			return err
		}

		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusOK,
			Body:   b,
		})
	}

	if len(path) == 3 && path[0] == "budgets" && path[2] == "accounts" {
		accounts, err := d.client.Accounts(ctx, path[1])
		if err != nil {
			return err
		}

		var selectableValues []SelectableValue

		for _, account := range accounts {
			selectableValues = append(selectableValues, SelectableValue{
				Label: account.Name,
				Value: account.ID,
			})
		}

		b, err := json.Marshal(selectableValues)
		if err != nil {
			return err
		}

		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusOK,
			Body:   b,
		})
	}

	return sender.Send(&backend.CallResourceResponse{
		Status: http.StatusNotFound,
	})
}

func transactionsToWideFrame(txs []ynab.Transaction, groups map[string]ynab.CategoryGroup, alignBy string, queryType string) (*data.Frame, error) {
	var dates []time.Time
	var accounts []string
	var payees []string
	var memos []string
	var amounts []float64
	var categories []string
	var categoryGroups []*string

	for _, tx := range txs {
		date, err := time.Parse("2006-01-02", tx.Date)
		if err != nil {
			backend.Logger.Error(err.Error())
			continue
		}

		switch queryType {
		case "spending":
			if tx.Amount < 0 {
				dates = append(dates, date)
				payees = append(payees, tx.PayeeName)
				memos = append(memos, tx.Memo)
				accounts = append(accounts, tx.AccountName)
				amounts = append(amounts, -float64(tx.Amount)/1000)
				categories = append(categories, tx.CategoryName)

				if g, ok := groups[tx.CategoryID]; ok {
					val := g.Name
					categoryGroups = append(categoryGroups, &val)
				} else {
					categoryGroups = append(categoryGroups, nil)
				}
			}
		case "income":
			if tx.Amount > 0 {
				dates = append(dates, date)
				payees = append(payees, tx.PayeeName)
				memos = append(memos, tx.Memo)
				accounts = append(accounts, tx.AccountName)
				amounts = append(amounts, float64(tx.Amount)/1000)
				categories = append(categories, tx.CategoryName)

				if g, ok := groups[tx.CategoryID]; ok {
					val := g.Name
					categoryGroups = append(categoryGroups, &val)
				} else {
					categoryGroups = append(categoryGroups, nil)
				}
			}
		}

	}

	frame := data.NewFrame("transactions",
		data.NewField("time", nil, dates),
		data.NewField("account", nil, accounts),
		data.NewField("payee", nil, payees),
		data.NewField("amount", nil, amounts),
		data.NewField("memo", nil, memos),
		data.NewField("category", nil, categories),
		data.NewField("category_group", nil, categoryGroups),
	)

	switch alignBy {
	case "category":
		return alignFrameByField(frame, "category"), nil
	case "category_group":
		return alignFrameByField(frame, "category_group"), nil
	case "account":
		return alignFrameByField(frame, "account"), nil
	case "payee":
		return alignFrameByField(frame, "payee"), nil
	}

	return frame, nil
}

func alignFrameByField(sourceFrame *data.Frame, fieldName string) *data.Frame {
	alignmentField, n := sourceFrame.FieldByName(fieldName)
	if n < 0 {
		return sourceFrame
	}
	amountField, n := sourceFrame.FieldByName("amount")
	if n < 0 {
		return sourceFrame
	}

	var uniqueValues []string
	for i := 0; i < alignmentField.Len(); i++ {
		switch val := alignmentField.At(i).(type) {
		case string:
			if !sliceContains(uniqueValues, val) {
				uniqueValues = append(uniqueValues, val)
			}
		case *string:
			if val != nil && !sliceContains(uniqueValues, *val) {
				uniqueValues = append(uniqueValues, *val)
			}
		}

	}

	alignedFields := make(map[string]*data.Field)
	for _, c := range uniqueValues {
		alignedFields[c] = data.NewField(c, nil, make([]*float64, sourceFrame.Rows()))
	}

	for i := 0; i < sourceFrame.Rows(); i++ {
		switch val := alignmentField.At(i).(type) {
		case string:
			field, ok := alignedFields[val]
			if ok {
				val := amountField.At(i).(float64)
				field.Set(i, &val)
			}
		case *string:
			if val != nil {
				field, ok := alignedFields[*val]
				if ok {
					val := amountField.At(i).(float64)
					field.Set(i, &val)
				}
			}
		}

	}

	var newFields []*data.Field
	for _, f := range alignedFields {
		newFields = append(newFields, f)
	}

	timeField, _ := sourceFrame.FieldByName("time")

	return data.NewFrame("transactions", append([]*data.Field{timeField}, newFields...)...)
}

func sliceContains(slice []string, find string) bool {
	for _, str := range slice {
		if str == find {
			return true
		}
	}
	return false
}
