package httpServer

import (
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"logserver/common"
	"logserver/etcd"
	"logserver/kafka"
	"tutu_task/handler"
)

func NewRouter() *fasthttprouter.Router {
	var router *fasthttprouter.Router
	router = fasthttprouter.New()
	router.POST("/log/inses", LogToKafka)
	router.GET("/log/list", LogJobList)
	router.POST("/log/putjob", LogPutJob)
	router.GET("/log/deljob", LogDelJob)
	router.GET("/log/delall", LogAllDelJob)
	router.POST("/log/delbulk", LogBulkDelJob)
	return router
}

func LogJobList(ctx *fasthttp.RequestCtx) {
	defer func() {
		if err := recover(); err != nil {
			handler.DoJSONWrite(ctx, 400, handler.GenerateResp("", -1, "failed"))
			return
		}
	}()

	var (
		resp map[string]map[string]interface{}
		jobs []*common.Jobs
		err  error
		res  []byte
	)
	if jobs, err = etcd.GjobMgr.ListLogJobs(); err != nil {
		goto ERR
	}

	if res, err = json.Marshal(jobs); err != nil {
		goto ERR
	}

	resp = handler.GenerateResp(string(res), 0, "success")
	handler.DoJSONWrite(ctx, 200, resp)
	return

ERR:
	resp = handler.GenerateResp(err.Error(), -1, "failed")
	handler.DoJSONWrite(ctx, 400, resp)
	return
}

func LogPutJob(ctx *fasthttp.RequestCtx) {
	defer func() {
		if err := recover(); err != nil {
			handler.DoJSONWrite(ctx, 400, handler.GenerateResp("", -1, "failed"))
			return
		}
	}()
	var (
		postRes    []byte
		err        error
		mapResults map[string]interface{}
		key        string
		jobName    string
		topic      string
		pipeline   string
		index_type string
		resp       map[string]map[string]interface{}
	)
	postRes = ctx.PostBody()
	if err = json.Unmarshal(postRes, &mapResults); err != nil {
		goto ERR
	}

	key = mapResults["key"].(string)
	jobName = mapResults["job_name"].(string)
	topic = mapResults["topic"].(string)
	if v, ok := mapResults["pipeline"]; ok {
		pipeline = v.(string)
	}
	if v, ok := mapResults["index_type"]; ok {
		index_type = v.(string)
	}
	//common.Jobs{}

	if _, err = etcd.GjobMgr.AddNewLogjob(key, common.Jobs{
		JobName:   jobName,
		Topic:     topic,
		IndexType: index_type,
		Pipeline:  pipeline,
	}); err != nil {
		goto ERR
	}
	resp = handler.GenerateResp("添加成功", 0, "success")
	handler.DoJSONWrite(ctx, 200, resp)
	return

ERR:
	resp = handler.GenerateResp(err.Error(), -1, "failed")
	handler.DoJSONWrite(ctx, 400, resp)
	return
}

func LogDelJob(ctx *fasthttp.RequestCtx) {
	defer func() {
		if err := recover(); err != nil {
			handler.DoJSONWrite(ctx, 400, handler.GenerateResp("", -1, "failed"))
			return
		}
	}()
	var (
		resp      map[string]map[string]interface{}
		getValues *fasthttp.Args
		key       []byte
		job       *common.Jobs
		err       error
		res       []byte
	)
	getValues = ctx.QueryArgs()
	key = getValues.Peek("key")
	if job, err = etcd.GjobMgr.DelLogJob(string(key)); err != nil {
		goto ERR
	}
	if res, err = json.Marshal(job); err == nil {
		resp = handler.GenerateResp(string(res), 0, "success")
		handler.DoJSONWrite(ctx, 200, resp)
		return
	} else {
		resp = handler.GenerateResp("删除成功", 0, "success")
		handler.DoJSONWrite(ctx, 200, resp)
		return
	}

ERR:
	resp = handler.GenerateResp(err.Error(), -1, "failed")
	handler.DoJSONWrite(ctx, 400, resp)
	return
}


//批量删除接口
func LogBulkDelJob(ctx *fasthttp.RequestCtx) {
	defer func() {
		if err := recover(); err != nil {
			handler.DoJSONWrite(ctx, 400, handler.GenerateResp("", -1, "failed"))
			return
		}
	}()
	var (
		postRes    []byte
		resp       map[string]map[string]interface{}
		mapResults map[string]interface{}

		keyStrs []string

		err error
	)
	postRes = ctx.PostBody()
	if err = json.Unmarshal(postRes, &mapResults); err != nil {
		goto ERR
	}
	keyStrs = make([]string, 0)
	if v, ok := mapResults["keys"]; ok {
		for _, sonv := range v.([]interface{}) {
			keyStrs = append(keyStrs, sonv.(string))
		}
	} else {
		goto ERR
	}
	if err = etcd.GjobMgr.BulkDelLogJob(keyStrs); err != nil {
		goto ERR
	}

	resp = handler.GenerateResp("删除成功", 0, "success")
	handler.DoJSONWrite(ctx, 200, resp)
	return

ERR:
	resp = handler.GenerateResp(err.Error(), -1, "failed")
	handler.DoJSONWrite(ctx, 400, resp)
	return
}

//全量删除接口
func LogAllDelJob(ctx *fasthttp.RequestCtx) {
	defer func() {
		if err := recover(); err != nil {
			handler.DoJSONWrite(ctx, 400, handler.GenerateResp("", -1, "failed"))
			return
		}
	}()
	var (
		resp       map[string]map[string]interface{}
		err error
	)

	if err = etcd.GjobMgr.DeleteAllJob(); err != nil{
		goto ERR
	}

	resp = handler.GenerateResp("删除成功", 0, "success")
	handler.DoJSONWrite(ctx, 200, resp)
	return

ERR:
	resp = handler.GenerateResp(err.Error(), -1, "failed")
	handler.DoJSONWrite(ctx, 400, resp)
	return
}


func LogToKafka(ctx *fasthttp.RequestCtx) {
	defer func() {
		if err := recover(); err != nil {
			handler.DoJSONWrite(ctx, 400, handler.GenerateResp("", -2, "failed"))
			return
		}
	}()
	var (
		postRes    []byte
		mapResults map[string]interface{}
		err        error
		topic      string
		logs       []interface{}
		kafkaLogs  []string
		resp       map[string]map[string]interface{}
	)
	postRes = ctx.PostBody()
	if err = json.Unmarshal(postRes, &mapResults); err != nil {
		goto ERR
	}
	topic = mapResults["topic"].(string)
	logs = mapResults["logs"].([]interface{})
	kafkaLogs = make([]string, 0)
	for _, v := range logs {
		kafkaLogs = append(kafkaLogs, v.(string))
	}
	fmt.Println(kafkaLogs)
	if err = kafka.SendToKafka(kafkaLogs, topic); err != nil {
		goto ERR
	}
	resp = handler.GenerateResp("插入成功", 0, "success")
	handler.DoJSONWrite(ctx, 200, resp)
	return
ERR:
	resp = handler.GenerateResp(err.Error(), -1, "failed")
	handler.DoJSONWrite(ctx, 400, resp)
	return
}

func InitHttpServer() error {
	var (
		router *fasthttprouter.Router
		err    error
	)
	router = NewRouter()
	if err = fasthttp.ListenAndServe(":8999", router.Handler); err != nil {
		fmt.Println("start fasthttp fail:", err.Error())
	}
	return err
}
