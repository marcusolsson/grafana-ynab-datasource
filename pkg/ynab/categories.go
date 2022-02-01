package ynab

import (
	"context"
	"fmt"
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
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/budgets/%s/categories", budgetID))
	if err != nil {
		return nil, err
	}

	var payload struct {
		Data struct {
			CategoryGroups []CategoryGroup `json:"category_groups"`
		} `json:"data"`
	}

	if err := c.do(req, &payload); err != nil {
		return nil, err
	}

	return payload.Data.CategoryGroups, nil
}
