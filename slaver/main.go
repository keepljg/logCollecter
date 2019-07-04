package logserver

import (
	"fmt"
	"logserver/logs"
	"logserver/slaver/conf"
	"logserver/slaver/es"
	"logserver/slaver/etcd"
	"logserver/slaver/httpServer"
	"logserver/slaver/kafka"
	"logserver/slaver/scheduler"
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
	if err = conf.InitConf(); err != nil {
		goto ERR
	}

	initEnv()
	//// 初始化logger
	//if err = serverlogs.InitLogger(); err != nil {
	//	goto ERR
	//}

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
	fmt.Println("123")
	if err = httpServer.InitHttpServer(); err != nil {
		goto ERR
	}

ERR:
	logs.ERROR(err)
	return
}
