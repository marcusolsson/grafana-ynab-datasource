package main

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/marcusolsson/grafana-ynab-datasource/pkg/ynab"
)

func TestAligner_Daily(t *testing.T) {
	for _, tt := range []struct {
		In  []ynab.Transaction
		Out []Measurement
	}{
		{
			In: []ynab.Transaction{
				{Date: "2006-01-02", Amount: -1000},
			},
			Out: []Measurement{
				{Time: time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC), Value: -1000},
			},
		},
		{
			In: []ynab.Transaction{
				{Date: "2006-01-02", Amount: 1000},
				{Date: "2006-01-02", Amount: 500},
				{Date: "2006-01-02", Amount: 250},
			},
			Out: []Measurement{
				{Time: time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC), Value: 250},
			},
		},
		{
			In: []ynab.Transaction{
				{Date: "2006-01-01", Amount: 1000},
				{Date: "2006-01-03", Amount: 500},
				{Date: "2006-01-05", Amount: 250},
			},
			Out: []Measurement{
				{Time: time.Date(2006, time.January, 1, 0, 0, 0, 0, time.UTC), Value: 1000},
				{Time: time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC), Value: 1000},
				{Time: time.Date(2006, time.January, 3, 0, 0, 0, 0, time.UTC), Value: 500},
				{Time: time.Date(2006, time.January, 4, 0, 0, 0, 0, time.UTC), Value: 500},
				{Time: time.Date(2006, time.January, 5, 0, 0, 0, 0, time.UTC), Value: 250},
			},
		},
	} {
		var alignable TimeSeriesTransactions = tt.In

		out, err := Regularize(alignable, PeriodDaily, alignLast)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(tt.Out, out); diff != "" {
			t.Error(diff)
		}
	}
}

func TestAligner_Monthly(t *testing.T) {
	for _, tt := range []struct {
		In  []ynab.Transaction
		Out []Measurement
	}{
		{
			In: []ynab.Transaction{
				{Date: "2006-01-02", Amount: -1000},
			},
			Out: []Measurement{
				{Time: time.Date(2006, time.January, 1, 0, 0, 0, 0, time.UTC), Value: -1000},
			},
		},
		{
			In: []ynab.Transaction{
				{Date: "2006-01-02", Amount: 1000},
				{Date: "2006-01-02", Amount: 500},
				{Date: "2006-01-02", Amount: 250},
				{Date: "2006-03-02", Amount: 500},
			},
			Out: []Measurement{
				{Time: time.Date(2006, time.January, 1, 0, 0, 0, 0, time.UTC), Value: 250},
				{Time: time.Date(2006, time.February, 1, 0, 0, 0, 0, time.UTC), Value: 250},
				{Time: time.Date(2006, time.March, 1, 0, 0, 0, 0, time.UTC), Value: 500},
			},
		},
	} {
		var alignable TimeSeriesTransactions = tt.In

		out, err := Regularize(alignable, PeriodMonthly, alignLast)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(tt.Out, out); diff != "" {
			t.Error(diff)
		}
	}
}
