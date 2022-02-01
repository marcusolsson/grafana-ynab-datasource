package ynab_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/marcusolsson/grafana-ynab-datasource/pkg/ynab"
)

func TestTransactions_200(t *testing.T) {
	var invoked bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		invoked = true

		if path := "/budgets/my-budget/transactions"; req.URL.Path != path {
			t.Errorf("unexpected path: want = %q; got = %q", path, req.URL.Path)
		}

		b, err := ioutil.ReadFile(filepath.Join("testdata", "transactions_200.json"))
		if err != nil {
			t.Fatal(err)
		}

		w.Write(b)
	}))

	client := ynab.NewClient("api-token")
	client.BaseURL = srv.URL

	ctx := context.Background()

	got, err := client.Transactions(ctx, "my-budget", "", "2006-01-02")
	if err != nil {
		t.Fatal(err)
	}

	if !invoked {
		t.Fatal("missing request")
	}

	want := []ynab.Transaction{
		{
			ID:                    "string",
			Date:                  "string",
			Memo:                  "string",
			Cleared:               "cleared",
			Approved:              true,
			FlagColor:             "red",
			AccountID:             "string",
			PayeeID:               "string",
			CategoryID:            "string",
			TransferAccountID:     "string",
			TransferTransactionID: "string",
			MatchedTransactionID:  "string",
			ImportID:              "string",
			Deleted:               true,
			AccountName:           "string",
			PayeeName:             "string",
			CategoryName:          "string",
			Subtransactions: []ynab.Transaction{
				{
					ID:                    "string",
					Memo:                  "string",
					PayeeID:               "string",
					CategoryID:            "string",
					TransferAccountID:     "string",
					TransferTransactionID: "string",
					Deleted:               true,
					PayeeName:             "string",
					CategoryName:          "string",
				},
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}
