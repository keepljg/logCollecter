package es

import (
	"context"
	"github.com/astaxie/beego/logs"
	"github.com/olivere/elastic"
	"logserver/configs"
)

func InitElasticClient() error {
	var (
		err error
		searchClient *ElasticClient
		pingRes *elastic.PingResult
		code int
	)

	searchClient = new(ElasticClient)
	searchClient.Client, err = elastic.NewClient(elastic.SetURL(configs.AppConfig.EsAddr))
	if err != nil {
		panic(err)
	}
	pingRes, code, err = searchClient.Client.Ping(configs.AppConfig.EsAddr).Do(context.Background())
	if err != nil {
		panic(err)
	}
	logs.Info("Elasticsearch returned with code %d and version %s\n", code, pingRes.Version.Number)

	GelasticCli = searchClient
	return err
}

//批量创建文档
func (e *ElasticClient) CreateBulkDocument(index string, docs []interface{}, pipeLine string) error {
	var(
		bulkService *elastic.BulkService
		err error
	)

	bulkService = e.Client.Bulk().Index(index).Type("doc")
	for i := 0; i < len(docs); i++ {
		bulkService.Add(elastic.NewBulkIndexRequest().Doc(docs[i])).Pipeline(pipeLine)
	}

	//Commit
	for trytimes := 10; trytimes > 0; trytimes-- {
		_, err = bulkService.Do(context.TODO())
		if err == nil {
			break
		}
	}
	if err == nil {
	} else {
		logs.Error(err)
		return err
	}

	return nil
}