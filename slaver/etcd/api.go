package etcd

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"go.etcd.io/etcd/clientv3"
	"logserver/slaver/common"
	"logserver/slaver/configs"
	"time"
)

func InitJobMgr() error {
	gin.Logger()
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
	)
	config = clientv3.Config{
		Endpoints:            []string{configs.AppConfig.EtcdAddr},
		DialKeepAliveTimeout: time.Duration(configs.AppConfig.EtcdTimeOut) * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		return err
	}
	GjobMgr = &JobMgr{
		Client:  client,
		Kv:      clientv3.NewKV(client),
		Lease:   clientv3.NewLease(client),
		Watcher: clientv3.NewWatcher(client),
	}
	return nil
}

// 新增log 任务
func (this *JobMgr) AddNewLogjob(key string, job common.Jobs) (*common.Jobs, error) {
	var (
		jobKey string
		err    error
		val    []byte
		putRes *clientv3.PutResponse
	)
	jobKey = configs.AppConfig.JobSave + key
	if val, err = json.Marshal(job); err != nil {
		logs.Error(err)
		return nil, err
	}
	if putRes, err = this.Kv.Put(context.TODO(), jobKey, string(val), clientv3.WithPrevKV()); err != nil {
		logs.Error(err)
		return nil, err
	}
	if putRes.PrevKv != nil {
		var oldJob common.Jobs
		err = json.Unmarshal([]byte(putRes.PrevKv.Value), &oldJob)
		if err != nil {
			return nil, nil
		} else {
			return &oldJob, nil
		}
	} else {
		return nil, nil
	}
}

// 删除log 任务
func (this *JobMgr) DelLogJob(key string) (*common.Jobs, error) {
	var (
		jobKey string
		err    error
		delRes *clientv3.DeleteResponse
	)
	jobKey = configs.AppConfig.JobSave + key
	if delRes, err = this.Kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		logs.Error(err)
		return nil, err
	}
	if len(delRes.PrevKvs) > 0 {
		var oldJob common.Jobs
		if err = json.Unmarshal([]byte(delRes.PrevKvs[0].Value), &oldJob); err != nil {
			return nil, nil
		} else {
			return &oldJob, nil
		}
	}
	return nil, nil
}

// 批量删除log 任务
func (this *JobMgr) BulkDelLogJob(keys []string) error {
	var (
		jobKey string
		err    error
	)
	for _, key := range keys {
		jobKey = configs.AppConfig.JobSave + key
		if _, err = this.Kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
			logs.Error(err)
			return err
		}
	}
	return nil
}

// 删除所有任务
func (this *JobMgr) DeleteAllJob() error {
	var (
		jobs    []*common.Jobs
		jobKeys []string
		err     error
	)
	if jobs, err = this.ListLogJobs(); err != nil {
		return err
	}
	for _, v := range jobs {
		jobKeys = append(jobKeys, v.Topic)
	}
	if err = this.BulkDelLogJob(jobKeys); err != nil {
		return err
	}
	return nil
}

// list 所有的log任务
func (this *JobMgr) ListLogJobs() ([]*common.Jobs, error) {
	var (
		err    error
		getRes *clientv3.GetResponse
		jobs   []*common.Jobs
	)
	if getRes, err = this.Kv.Get(context.TODO(), configs.AppConfig.JobSave, clientv3.WithPrefix()); err != nil {
		logs.Error(err)
		return nil, err
	}
	jobs = make([]*common.Jobs, 0)
	for _, v := range getRes.Kvs {
		var job common.Jobs
		if err = json.Unmarshal(v.Value, &job); err == nil {
			jobs = append(jobs, &job)
		}
	}
	return jobs, nil
}

// list 所有的抢到锁到任务
func (this *JobMgr) ListLogLocks() (map[string]string, error) {
	var (
		err    error
		getRes *clientv3.GetResponse
		locks  map[string]string
	)
	if getRes, err = this.Kv.Get(context.TODO(), configs.AppConfig.JobLock, clientv3.WithPrefix()); err != nil {
		logs.Error(err)
		return nil, err
	}
	locks = make(map[string]string)
	for _, v := range getRes.Kvs {
		locks[common.ExtractLockName(string(v.Key))] = ""
	}
	return locks, nil
}

// 创建事物锁
func (this *JobMgr) CreateJobLock(jobName string) *JobLock {
	return InitJobLock(jobName, this.Kv, this.Lease)
}
