package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/marcusolsson/grafana-ynab-datasource/pkg/ynab"
)

// DataSourceQuery represents the query that's constructed by the query editor
// and sent by the dashboard.
type DataSourceQuery struct {
	BudgetID          string   `json:"budgetId"`
	AccountIDs        []string `json:"accountIds"`
	AlignBy           string   `json:"alignBy"`
	Period            string   `json:"period"`
	QueryType         string   `json:"queryType"`
	TransactionFilter string   `json:"transactionFilter"`
	SpendingFilter    string   `json:"spendingFilter"`
	NetWorthFilter    string   `json:"netWorthFilter"`
}

type YNABDataSource struct {
	client *ynab.CacheClient
}

// NewYNABDatasource creates a new datasource instance. This function is called
// any time the settings are updated.
func NewYNABDatasource(client *ynab.CacheClient) *YNABDataSource {
	return &YNABDataSource{client: client}
}

func (d *YNABDataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()

	for _, query := range req.Queries {
		response.Responses[query.RefID] = d.query(ctx, req.PluginContext, query)
	}

	return response, nil
}

func (d *YNABDataSource) query(ctx context.Context, _ backend.PluginContext, dataQuery backend.DataQuery) backend.DataResponse {
	var dsQuery DataSourceQuery
	if err := json.Unmarshal(dataQuery.JSON, &dsQuery); err != nil {
		return backend.DataResponse{Error: err}
	}

	if dsQuery.BudgetID == "" {
		return backend.DataResponse{}
	}

	var (
		from = dataQuery.TimeRange.From
		to   = dataQuery.TimeRange.To
	)

	switch dsQuery.QueryType {
	case "transactions":
		return d.queryTransactions(ctx, dsQuery, from, to)
	case "net_worth":
		return d.queryNetWorth(ctx, dsQuery, from, to)
	case "spending":
		return d.querySpending(ctx, dsQuery, from, to)
	default:
		return backend.DataResponse{
			Error: errors.New("unsupported query type"),
		}
	}
}

// CheckHealth runs when the user presses "Save & Test" in the data source settings. Tests that the user can access
// their workspace.
func (d *YNABDataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if err := d.client.Test(ctx); err != nil {
		switch e := err.(type) {
		case ynab.APIError:
			return &backend.CheckHealthResult{
				Status:  backend.HealthStatusError,
				Message: e.Detail,
			}, nil
		default:
			backend.Logger.Error("health check failed", "err", e)

			return &backend.CheckHealthResult{
				Status:  backend.HealthStatusError,
				Message: "An error occurred. Check server logs for more details.",
			}, nil
		}
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Successfully connected to YNAB.",
	}, nil
}

// CallResources exposes a REST API with support operations for the query editor.
func (d *YNABDataSource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	if req.Path == "budgets" {
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

	return sender.Send(&backend.CallResourceResponse{
		Status: http.StatusNotFound,
	})
}
