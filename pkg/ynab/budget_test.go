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

func TestBudget_200(t *testing.T) {
	var invoked bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		invoked = true

		if path := "/budgets"; req.URL.Path != path {
			t.Errorf("unexpected path: want = %q; got = %q", path, req.URL.Path)
		}

		b, err := ioutil.ReadFile(filepath.Join("testdata", "budgets_200.json"))
		if err != nil {
			t.Fatal(err)
		}

		w.Write(b)
	}))

	client := ynab.NewClient("api-token")
	client.BaseURL = srv.URL

	ctx := context.Background()

	got, err := client.Budgets(ctx, true)
	if err != nil {
		t.Fatal(err)
	}

	if !invoked {
		t.Fatal("missing request")
	}

	want := []ynab.Budget{
		{
			ID:   "string",
			Name: "string",
			CurrencyFormat: ynab.CurrencyFormat{
				ISOCode: "string",
			},
			Accounts: []ynab.Account{
				{
					ID:                  "string",
					Name:                "string",
					Type:                "checking",
					OnBudget:            true,
					Closed:              true,
					Note:                "string",
					TransferPayeeID:     "string",
					DirectImportLinked:  true,
					DirectImportInError: true,
					Deleted:             true,
				},
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}

func TestBudget_404(t *testing.T) {
	var invoked bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		invoked = true

		if path := "/budgets"; req.URL.Path != path {
			t.Errorf("unexpected path: want = %q; got = %q", path, req.URL.Path)
		}

		b, err := ioutil.ReadFile(filepath.Join("testdata", "api_error.json"))
		if err != nil {
			t.Fatal(err)
		}

		w.WriteHeader(404)

		w.Write(b)
	}))

	client := ynab.NewClient("api-token")
	client.BaseURL = srv.URL

	ctx := context.Background()

	_, got := client.Budgets(ctx, true)
	if got == nil {
		t.Fatal("missing error")
	}

	if !invoked {
		t.Fatal("missing request")
	}

	want := ynab.APIError{
		ID:     "string",
		Name:   "string",
		Detail: "string",
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}
