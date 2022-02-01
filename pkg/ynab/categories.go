package ynab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type CategoryGroup struct {
	ID         string     `json:"id,omitempty"`
	Name       string     `json:"name,omitempty"`
	Hidden     bool       `json:"hidden,omitempty"`
	Deleted    bool       `json:"deleted,omitempty"`
	Categories []Category `json:"categories,omitempty"`
}

type Category struct {
	ID                      string `json:"id,omitempty"`
	CategoryGroupID         string `json:"category_group_id,omitempty"`
	Name                    string `json:"name,omitempty"`
	OriginalCategoryGroupID string `json:"original_category_group_id,omitempty"`
	Budgeted                int64  `json:"budgeted,omitempty"`
	Activity                int64  `json:"activity,omitempty"`
	Balance                 int64  `json:"balance,omitempty"`
}

func (c *Client) Categories(ctx context.Context, budgetID string) ([]CategoryGroup, error) {
	backend.Logger.Error("YNABClient.Categories()", "budgetID", budgetID)

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+fmt.Sprintf("/budgets/%s/categories", budgetID), nil)
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
			CategoryGroups []CategoryGroup `json:"category_groups"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Data.CategoryGroups, nil
}
