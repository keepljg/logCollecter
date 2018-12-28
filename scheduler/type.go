package scheduler

import "logserver/common"

type Scheduler struct {
	jobEventChan chan *common.JobEvent
	jobWorkTable map[string]*common.JobWorkInfo
}

var Gscheduler *Scheduler
