package common

import "context"

// log 任务
type Jobs struct {
	JobName   string `json:"job_name"`
	Topic     string `json:"topic"`
	IndexType string `json:"index_type"`
	Pipeline  string `json:"pipeline"`
}

// log 任务事件
type JobEvent struct {
	EventType int
	Job       *Jobs
}

type JobWorkInfo struct {
	Job        *Jobs
	ConText    context.Context
	CancelFunc context.CancelFunc
}
