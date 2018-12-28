package etcd

import (
	"context"
	"go.etcd.io/etcd/clientv3"
)

var GjobMgr *JobMgr

// 任务管理器
type JobMgr struct {
	Client  *clientv3.Client
	Kv      clientv3.KV
	Lease   clientv3.Lease   // 租约
	Watcher clientv3.Watcher // 监听
}

// 分布式锁
type JobLock struct {
	// etcd客户端
	kv         clientv3.KV
	lease      clientv3.Lease     // 租约
	jobName    string             // 任务名
	cancelFunc context.CancelFunc // 用于终止自动续租
	leaseId    clientv3.LeaseID   // 租约ID
	isLocked   bool               // 是否上锁成功
}
