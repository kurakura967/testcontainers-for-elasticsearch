package testing

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/elasticsearch"
)

type ElasticsearchContainer struct {
	*elasticsearch.ElasticsearchContainer
}

func SetupElasticsearch(t *testing.T) *ElasticsearchContainer {
	t.Helper()
	ctx := context.Background()

	container, err := elasticsearch.RunContainer(ctx,
		testcontainers.WithImage("docker.elastic.co/elasticsearch/elasticsearch:8.7.1"),
		elasticsearch.WithPassword("PASSWORD"),
	)
	if err != nil {
		log.Fatalf("failed to start elasticsearch container: %v", err)
	}

	return &ElasticsearchContainer{
		container,
	}
}

func (container *ElasticsearchContainer) Down() {
	ctx := context.Background()
	if err := container.Terminate(ctx); err != nil {
		log.Fatalf("failed to stop elasticsearch container: %v", err)
	}
}

func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	bt, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	return bt
}

type esResponse struct {
	Hits esHitsOuter `json:"hits"`
}

type esHitsOuter struct {
	Hits []esHitsInner `json:"hits"`
}

type esHitsInner struct {
	Source esSource `json:"_source"`
}

type esSource struct {
	Tile   string `json:"tile"`
	Artist string `json:"artist"`
}

func AssertSearchResponse(t *testing.T, want, got []byte) {
	t.Helper()

	var gotRes, wantRes esResponse

	if err := json.Unmarshal(got, &gotRes); err != nil {
		t.Fatalf("failed to unmarshal JSON data: %v", err)
	}

	if err := json.Unmarshal(want, &wantRes); err != nil {
		t.Fatalf("failed to unmarshal JSON data: %v", err)
	}

	if diff := cmp.Diff(wantRes, gotRes); diff != "" {
		t.Errorf("unexpected response (-want, +got): %s\n", diff)
	}
}
