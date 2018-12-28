package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"logserver/common"
	"testing"
	"time"
)

func TestInitJobMgr(t *testing.T) {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
	)
	config = clientv3.Config{
		Endpoints:            []string{"23.236.115.227:2379"},
		DialKeepAliveTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		panic(err)
	}
	jobMgr := &JobMgr{
		Client:  client,
		Kv:      clientv3.NewKV(client),
		Lease:   clientv3.NewLease(client),
		Watcher: clientv3.NewWatcher(client),
	}
	resp,_ := jobMgr.Kv.Get(context.TODO(), common.JOB_LOCK_DIR + "tutuapp_test")
	for _, v := range resp.Kvs {
		fmt.Println(string(v.Value))
	}
	//resp := clientv3.OpGet(common.JOB_LOCK_DIR + "tutuapp_test",clientv3.WithPrefix())
	//fmt.Println(string(resp.KeyBytes()))
	//fmt.Println(string(resp.ValueBytes()))
}