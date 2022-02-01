package ynab_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/marcusolsson/grafana-ynab-datasource/pkg/ynab"
)

func TestBalance(t *testing.T) {
	for _, tt := range []struct {
		In  []ynab.Transaction
		Out []ynab.Balance
	}{
		{
			In: []ynab.Transaction{
				{Date: "2006-01-02", Amount: -1000},
			},
			Out: []ynab.Balance{
				{Date: "2006-01-02", Amount: -1000},
			},
		},
		{
			In: []ynab.Transaction{
				{Date: "2006-01-02", Amount: -1000},
				{Date: "2006-01-02", Amount: 1000},
			},
			Out: []ynab.Balance{
				{Date: "2006-01-02", Amount: -1000},
				{Date: "2006-01-02", Amount: 0},
			},
		},
		{
			In: []ynab.Transaction{
				{Date: "2006-01-01", Amount: 500},
				{Date: "2006-01-02", Amount: -1000},
				{Date: "2006-01-02", Amount: 1000},
				{Date: "2006-01-03", Amount: 1000},
			},
			Out: []ynab.Balance{
				{Date: "2006-01-01", Amount: 500},
				{Date: "2006-01-02", Amount: -500},
				{Date: "2006-01-02", Amount: 500},
				{Date: "2006-01-03", Amount: 1500},
			},
		},
		{
			In: []ynab.Transaction{
				{Date: "2006-01-01", Amount: 500},
				{Date: "2006-01-03", Amount: -1000},
			},
			Out: []ynab.Balance{
				{Date: "2006-01-01", Amount: 500},
				{Date: "2006-01-03", Amount: -500},
			},
		},
	} {
		out, err := ynab.RunningBalance(tt.In)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(tt.Out, out); diff != "" {
			t.Error(diff)
		}
	}
}
