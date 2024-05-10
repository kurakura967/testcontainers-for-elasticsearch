package testing

import (
	"context"

	repository "github.com/kurakura967/testcontainers-for-elasticsearch/elasticsearch"
	"github.com/kurakura967/testcontainers-for-elasticsearch/testing/data"
)

func SetupIndex(repo repository.Elasticsearch, index, mappings string) {
	ctx := context.Background()
	found, err := repo.IsExistsIndex(ctx, index)
	if err != nil {
		panic(err)
	}
	// 既に存在する場合は削除してから作成
	if found {
		err := repo.DeleteIndex(ctx, index)
		if err != nil {
			panic(err)
		}
	}

	err = repo.CreateIndex(ctx, index, mappings)
	if err != nil {
		panic(err)
	}

	for _, doc := range data.DummyDocuments {
		err := repo.InsertDocument(ctx, index, doc)
		if err != nil {
			panic(err)
		}
	}
}
