package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/kurakura967/testcontainers-for-elasticsearch/model"
)

type esClient struct {
	client *elasticsearch.Client
}

func NewElasticsearch(cfg elasticsearch.Config) (Elasticsearch, error) {
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &esClient{client}, nil
}

func (e *esClient) Search(ctx context.Context, index, query string) ([]byte, error) {
	reqCtx, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(reqCtx, e.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// TODO: Errorレスポンスの構造体を定義する
	// see: https://github.com/olivere/elastic/blob/release-branch.v7/errors.go
	if res.IsError() {
		var e errResponse
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("faild to read err response body: %w", err)
		}
		if err := json.Unmarshal(body, &e); err != nil {
			return nil, fmt.Errorf("faild to unmarsal err response body: %w", err)
		}
		return nil, fmt.Errorf("failt to search: [%d] %s", e.Status, e.Error.Cause[0].Reason)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (e *esClient) CreateIndex(ctx context.Context, index, mapping string) error {
	reqCtx, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	req := esapi.IndicesCreateRequest{
		Index: index,
		Body:  strings.NewReader(mapping),
	}

	res, err := req.Do(reqCtx, e.client)
	if err != nil {
		return err
	}

	if res.IsError() {
		return fmt.Errorf("faild to create index: %s", res.Status())
	}

	return nil
}

func (e *esClient) DeleteIndex(ctx context.Context, index string) error {
	reqCtx, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	req := esapi.IndicesDeleteRequest{
		Index: []string{index},
	}

	res, err := req.Do(reqCtx, e.client)
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("faild to delete index: %s", res.Status())
	}
	return nil
}

func (e *esClient) IsExistsIndex(ctx context.Context, index string) (bool, error) {
	reqCtx, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	req := esapi.IndicesExistsRequest{
		Index: []string{index},
	}

	res, err := req.Do(reqCtx, e.client)
	if err != nil {
		return false, err
	}
	s := res.StatusCode
	switch s {
	case http.StatusOK:
		log.Printf("index %s is already exits", index)
		return true, nil
	case http.StatusNotFound:
		log.Printf("index %s not found", index)
		return false, nil
	default:
		return false, nil
	}
}

func (e *esClient) InsertDocument(ctx context.Context, index string, document model.Document) error {
	reqCtx, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	data, err := json.Marshal(document)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: document.ID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(reqCtx, e.client)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}

type errResponse struct {
	Status int      `json:"status"`
	Error  errCause `json:"error"`
}

type errCause struct {
	Cause []errReason `json:"root_cause"`
}

type errReason struct {
	Reason string `json:"reason"`
}
