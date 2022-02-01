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

func TestCategories_200(t *testing.T) {
	var invoked bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		invoked = true

		if path := "/budgets/my-budgets/categories"; req.URL.Path != path {
			t.Errorf("unexpected path: want = %q; got = %q", path, req.URL.Path)
		}

		b, err := ioutil.ReadFile(filepath.Join("testdata", "categories_200.json"))
		if err != nil {
			t.Fatal(err)
		}

		w.Write(b)
	}))

	client := ynab.NewClient("api-token")
	client.BaseURL = srv.URL

	ctx := context.Background()

	got, err := client.Categories(ctx, "my-budgets")
	if err != nil {
		t.Fatal(err)
	}

	if !invoked {
		t.Fatal("missing request")
	}

	want := []ynab.CategoryGroup{
		{
			ID:      "string",
			Name:    "string",
			Hidden:  true,
			Deleted: true,
			Categories: []ynab.Category{
				{
					ID:                      "string",
					CategoryGroupID:         "string",
					Name:                    "string",
					OriginalCategoryGroupID: "string",
					Budgeted:                0,
					Activity:                0,
					Balance:                 0,
				},
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}
