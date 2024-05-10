package repository_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	repository "github.com/kurakura967/testcontainers-for-elasticsearch/elasticsearch"
	testingHelper "github.com/kurakura967/testcontainers-for-elasticsearch/testing"
)

const defaultTestIndex = "test_index"

func TestSearch(t *testing.T) {
	container := testingHelper.SetupElasticsearch(t)
	t.Cleanup(container.Down)

	config := elasticsearch.Config{
		Addresses:    []string{container.Settings.Address},
		Username:     container.Settings.Username,
		Password:     container.Settings.Password,
		DisableRetry: false,
		CACert:       container.Settings.CACert,
	}

	client, err := repository.NewElasticsearch(config)
	if err != nil {
		t.Fatal(err)
	}

	testcases := []struct {
		name                 string
		index                string
		queryFileName        string
		wantResponseFileName string
		wantErr              bool
	}{
		{
			"正常系(200)",
			defaultTestIndex,
			filepath.Join("testdata", "ok_query.json"),
			filepath.Join("testdata", "ok_response.json"),
			false,
		},
		{
			"異常系(400-BadRequest)",
			defaultTestIndex,
			filepath.Join("testdata", "bad_request_query.json"),
			"",
			true,
		},
		{
			"異常系(404-IndexNotFound)",
			"not_found_index",
			filepath.Join("testdata", "ok_query.json"),
			"",
			true,
		},
	}

	for _, testcase := range testcases {
		ctx := context.Background()
		mappings := testingHelper.LoadFile(t, filepath.Join("testdata", "mappings.json"))
		testingHelper.SetupIndex(client, defaultTestIndex, string(mappings))

		t.Run(testcase.name, func(t *testing.T) {
			query := testingHelper.LoadFile(t, testcase.queryFileName)
			got, err := client.Search(ctx, testcase.index, string(query))
			if (err != nil) != testcase.wantErr {
				t.Errorf("error: got %v, wantErr %v", err, testcase.wantErr)
				return
			}
			if testcase.wantResponseFileName != "" {
				want := testingHelper.LoadFile(t, testcase.wantResponseFileName)
				testingHelper.AssertSearchResponse(t, want, got)
			}
		})
	}
}
