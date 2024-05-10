package repository

import (
	"context"

	"github.com/kurakura967/testcontainers-for-elasticsearch/model"
)

type Elasticsearch interface {
	Search(ctx context.Context, index, query string) ([]byte, error)
	CreateIndex(ctx context.Context, index, mapping string) error
	IsExistsIndex(ctx context.Context, index string) (bool, error)
	DeleteIndex(ctx context.Context, index string) error
	InsertDocument(ctx context.Context, index string, document model.Document) error
}
