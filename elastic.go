package elastic_client

import (
	elasticv7 "github.com/olivere/elastic/v7"
)

type SearchRequest struct {
	Id    int
	Type  int
	Title string
	Body  string
}

type elastic struct {
	Connection *elasticv7.Client
	Url        string
}

func New() *elastic {
	return &elastic{}
}
