package scheduler

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"logserver/logs"
	"logserver/slaver/common"
	"logserver/slaver/conf"
	"logserver/slaver/etcd"
	"logserver/slaver/kafka"
	"time"
)

func InitScheduler() {
	NewScheduler()
	go Gscheduler.ScheuleLoop()
	Gscheduler.restartLogJob()
	go Gscheduler.CalculatingPressure()
	WatcherJobs()
}

func NewScheduler() {
	Gscheduler = &Scheduler{
		logCount:     make(chan int, 1000),
		JobEventChan: make(chan *common.JobEvent, 1000),
		JobWorkTable: make(map[string]*common.JobWorkInfo),
	}
}

// push log任务事件
func (this *Scheduler) PushJobEvent(event *common.JobEvent) {
	this.JobEventChan <- event
}

func (this *Scheduler) ScheuleLoop() {
	var (
		jobEvent *common.JobEvent
	)
	for {
		select {
		case jobEvent = <-this.JobEventChan:
			this.handleJobEvent(jobEvent)
		}
	}
}

//开启topic 的kafka消费
func (this *Scheduler) eventWorker(job *common.Jobs) {
	var (
		jobWorkInfo *common.JobWorkInfo
		jobLock     *etcd.JobLock
		err         error
	)
	jobLock = etcd.GjobMgr.CreateJobLock(job.Topic)
	err = jobLock.TryToLock()
	//defer jobLock.Unlock()
	if err == nil {
		jobWorkInfo = common.NewJobWorkInfo(job)
		if jobWork, ok := this.JobWorkTable[job.Topic]; !ok {
			this.JobWorkTable[job.Topic] = jobWorkInfo
			kafka.ConsumerFromKafka4(jobWorkInfo, jobLock, this.logCount)
		} else {
			// 重新开启新任务
			this.reEventWork(jobWork, job, jobLock)
		}
	}
}

// 更新任务
func (this *Scheduler) reEventWork(jobWork *common.JobWorkInfo, newJob *common.Jobs, lock *etcd.JobLock) {
	var (
		jobWorkInfo *common.JobWorkInfo
	)
	// 先关闭当前任务
	jobWork.CancelFunc()
	delete(this.JobWorkTable, jobWork.Job.Topic)
	jobWorkInfo = common.NewJobWorkInfo(newJob)
	this.JobWorkTable[newJob.Topic] = jobWorkInfo
	kafka.ConsumerFromKafka4(jobWorkInfo, lock, this.logCount)
}

// 处理日志任务
func (this *Scheduler) handleJobEvent(event *common.JobEvent) {
	switch event.EventType {
	case common.JOB_EVENT_SAVE:
		this.eventWorker(event.Job)
	case common.JOB_EVENT_DELETE:
		fmt.Println(this.JobWorkTable)
		fmt.Println(event.Job.Topic)
		if jobWork, ok := this.JobWorkTable[event.Job.Topic]; ok {
			jobWork.CancelFunc()
			delete(this.JobWorkTable, jobWork.Job.Topic)
		}
	}
}

// 进行监听log任务
func WatcherJobs() {
	var (
		getResp          *clientv3.GetResponse
		err              error
		job              *common.Jobs
		jobEvent         *common.JobEvent
		watcherReversion int64
		watchChan        clientv3.WatchChan
	)
	if getResp, err = etcd.GjobMgr.Kv.Get(context.TODO(), conf.JobConf.JobSave, clientv3.WithPrefix()); err != nil {
		logs.ERROR(err)
		return
	}
	fmt.Println( getResp, err)
	// 启动时先将etcd中的topic消费
	for _, v := range getResp.Kvs {
		if job, err = common.UnPackJob(v.Value); err == nil {
			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			Gscheduler.PushJobEvent(jobEvent)
		}
	}
	go func() {
		// 从getResp header 下一个版本进行监听
		watcherReversion = getResp.Header.Revision + 1
		watchChan = etcd.GjobMgr.Watcher.Watch(context.TODO(), conf.JobConf.JobSave, clientv3.WithRev(watcherReversion), clientv3.WithPrefix())
		for eachChan := range watchChan {
			for _, v := range eachChan.Events {
				switch v.Type {
				// 新的任务
				case clientv3.EventTypePut:
					if job, err = common.UnPackJob(v.Kv.Value); err == nil {
						jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
						Gscheduler.PushJobEvent(jobEvent)
					}
					// 删除任务
				case clientv3.EventTypeDelete:
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, &common.Jobs{
						Topic: common.ExtractJobName(string(v.Kv.Key)),
					})
					Gscheduler.PushJobEvent(jobEvent)
				}
			}
		}
	}()
}

// 容灾处理
func (this *Scheduler) restartLogJob() {
	var (
		t *time.Timer
	)
	t = time.NewTimer(time.Second * 60)
	go func() {
		for {
			select {
			case <-t.C:
				var (
					jobs  []*common.Jobs
					locks map[string]string
					err   error
				)
				// 所有在etcd中的log任务
				if jobs, err = etcd.GjobMgr.ListLogJobs(); err != nil {
					logs.ERROR(err)
					return
				}
				// 所有抢到锁的任务
				if locks, err = etcd.GjobMgr.ListLogLocks(); err != nil {
					logs.ERROR(err)
					return
				}
				for _, job := range jobs {
					if _, ok := locks[job.Topic]; !ok {
						this.PushJobEvent(&common.JobEvent{
							EventType: common.JOB_EVENT_SAVE,
							Job:       job,
						})
					}
				}
				t.Reset(time.Second * 60)
			}
		}
	}()
}

// 计算压力
func (this *Scheduler) CalculatingPressure() {
	var (
		jobCountSum, logCountSum int
	)
	counts := make([]int, 0)
	t := time.NewTimer(time.Second * 10)
	for {
		select {
		case count := <-this.logCount:
			counts = append(counts, count)
		case <-t.C:
			for _, count := range counts {
				logCountSum += count
			}
			logs.DEBUG("logCountSum is ", logCountSum)
			counts = make([]int, 0)
			// todo 上报信息

			jobCountSum = len(this.JobWorkTable)
			// todo 上报信息
			logs.INFO(jobCountSum)
			t.Reset(time.Second * 10)
		}
	}
	return
}


