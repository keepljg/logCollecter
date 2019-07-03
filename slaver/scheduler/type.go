package scheduler

import "logserver/slaver/common"

type Scheduler struct {
	JobEventChan chan *common.JobEvent
	JobWorkTable map[string]*common.JobWorkInfo
}

var Gscheduler *Scheduler
