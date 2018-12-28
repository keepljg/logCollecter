package es

import (
	"github.com/olivere/elastic"
)

type ElasticClient struct {
	Client *elastic.Client
}

var GelasticCli *ElasticClient
