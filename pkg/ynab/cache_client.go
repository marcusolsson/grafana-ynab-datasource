package ynab

import (
	"context"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type CacheClient struct {
	*Client

	cache *cache.Cache
}

func NewCacheClient(client *Client) *CacheClient {
	return &CacheClient{
		Client: client,
		cache:  cache.New(10*time.Minute, 60*time.Minute),
	}
}

func (c *CacheClient) Budgets(ctx context.Context, includeAccounts bool) ([]Budget, error) {
	budgets, found := c.cache.Get("budgets")
	if !found {
		budgets, err := c.Client.Budgets(ctx, includeAccounts)
		if err != nil {
			return nil, err
		}

		c.cache.Set("budgets", budgets, cache.DefaultExpiration)

		return budgets, nil
	}

	return budgets.([]Budget), nil
}

func (c *CacheClient) Categories(ctx context.Context, budgetID string) ([]CategoryGroup, error) {
	key := "categories/" + budgetID

	categories, found := c.cache.Get(key)
	if !found {
		categories, err := c.Client.Categories(ctx, budgetID)
		if err != nil {
			return nil, err
		}

		c.cache.Set(key, categories, cache.DefaultExpiration)

		return categories, nil
	}

	return categories.([]CategoryGroup), nil
}

func (c *CacheClient) Transactions(ctx context.Context, budgetID, accountID, sinceDate string) ([]Transaction, error) {
	key := fmt.Sprintf("transactions/%s/%s/%s", budgetID, accountID, sinceDate)

	txs, found := c.cache.Get(key)
	if !found {
		txs, err := c.Client.Transactions(ctx, budgetID, accountID, sinceDate)
		if err != nil {
			return nil, err
		}

		c.cache.Set(key, txs, cache.DefaultExpiration)

		return txs, nil
	}

	return txs.([]Transaction), nil
}
