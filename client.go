package elastic_client

import (
	"context"
	"errors"
	"reflect"

	elasticv7 "github.com/olivere/elastic/v7"
)

func (r elastic) CreateIndex(ctx context.Context, indexName string) bool {
	exists, err := r.Connection.IndexExists(indexName).Do(ctx)
	if err != nil {
		return false
	}

	if !exists {
		createIndex, err := r.Connection.CreateIndex(indexName).Do(ctx)
		if err != nil {
			return false
		}
		if !createIndex.Acknowledged {
			return false
		}
	}

	return true
}

func (r elastic) reDial(ctx context.Context) bool {
	_, _, err := r.Connection.Ping(r.Url).Do(ctx)

	return err == nil
}

func (r elastic) Add(indexName string, data interface{}) error {
	ctx := context.Background()
	if !r.reDial(ctx) {
		return errors.New("no redial")
	}

	r.CreateIndex(ctx, indexName)

	_, err := r.Connection.Index().
		Index(indexName).
		Type("page").
		BodyJson(data).
		Do(ctx)

	if err != nil {
		return err
	}

	_, err = r.Connection.Flush().Index(indexName).Do(ctx)
	return err
}

func (r elastic) Search(indexName string, query string) ([]SearchRequest, error) {
	result := []SearchRequest{}
	ctx := context.Background()
	if !r.reDial(ctx) {
		return nil, errors.New("no redial")
	}

	r.CreateIndex(ctx, indexName)

	termQuery := elasticv7.NewQueryStringQuery(query)
	searchResult, err := r.Connection.Search().
		Index(indexName).
		Query(termQuery).
		//Sort("user", true). // sort by "user" field, ascending
		From(0).Size(100).
		Pretty(true).
		Do(ctx)

	var ttyp SearchRequest
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(SearchRequest); ok {
			result = append(result, t)
		}
	}

	return result, err
}
