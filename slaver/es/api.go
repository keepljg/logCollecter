package es

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/olivere/elastic"
	"logserver/slaver/conf"
)

func InitElasticClient() error {
	var (
		err          error
		searchClient *ElasticClient
		//pingRes *elastic.PingResult
		//code int
	)

	searchClient = new(ElasticClient)
	searchClient.Client, err = elastic.NewClient(elastic.SetURL(conf.EsConf.Addr...))
	fmt.Println(conf.EsConf.Addr)
	if err != nil {
		panic(err)
	}

	GelasticCli = searchClient
	return err
}

//批量创建文档
func (e *ElasticClient) CreateBulkDocument(index string, docs []interface{}, pipeLine string) error {
	var (
		bulkService *elastic.BulkService
		err         error
	)

	bulkService = e.Client.Bulk().Index(index).Type("doc")
	for i := 0; i < len(docs); i++ {
		bulkService.Add(elastic.NewBulkIndexRequest().Doc(docs[i]).Pipeline(pipeLine))
	}

	//Commit
	for trytimes := 10; trytimes > 0; trytimes-- {
		_, err = bulkService.Do(context.TODO())
		if err == nil {
			break
		}
	}
	if err == nil {
		logs.Info(index + " insert data success")
	} else {
		logs.Error(err)
		return err
	}

	return nil
}

// 创建单个文档
func (e *ElasticClient) CreateSignDocument(index string, doc interface{}, pipeLine string) error {
	var err error
	_, err = e.Client.Index().Index(index).Type("doc").BodyJson(doc).Do(context.Background())
	logs.Error(err)
	return err
}
