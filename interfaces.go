package elastic_client

import (
	"context"
)

type Elastic interface {
	CreateIndex(ctx context.Context, indexName string) bool
	reDial(ctx context.Context) bool
	Add(indexName string, data interface{}) error
	Search(indexName string, query string) ([]SearchRequest, error)
}
