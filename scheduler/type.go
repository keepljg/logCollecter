package scheduler

import "logserver/common"

type Scheduler struct {
	JobEventChan chan *common.JobEvent
	JobWorkTable map[string]*common.JobWorkInfo
}

var Gscheduler *Scheduler
