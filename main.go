package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"logserver/configs"
	"logserver/es"
	"logserver/etcd"
	"logserver/httpServer"
	"logserver/kafka"
	"logserver/scheduler"
	"logserver/serverlogs"
	"runtime"
)

// 初始化线程数
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)
	// 初始化配置
	if err = configs.InitConf("ini", "./logserver.ini"); err != nil {
		goto ERR
	}

	initEnv()
	// 初始化logger
	if err = serverlogs.InitLogger(); err != nil {
		goto ERR
	}

	// 初始化etcd
	if err = etcd.InitJobMgr(); err != nil {
		goto ERR
	}

	//
	if err = kafka.InitKafka(); err != nil {
		goto ERR
	}

	if err = es.InitElasticClient(); err != nil {
		goto ERR
	}

	scheduler.InitScheduler()

	if err = httpServer.InitHttpServer(); err != nil {
		goto ERR
	}

ERR:
	fmt.Println(err)
	logs.Error(err)
	return
}
