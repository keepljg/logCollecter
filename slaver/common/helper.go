package common

import (
	"context"
	"encoding/json"
	"logserver/slaver/conf"
	"strings"
	"time"
)

func BuildJobEvent(eventType int, job *Jobs) *JobEvent {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

func UnPackJob(value []byte) (*Jobs, error) {
	var (
		err error
		job Jobs
	)
	if err = json.Unmarshal(value, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

func NewJobWorkInfo(job *Jobs) *JobWorkInfo {
	var (
		ctx        context.Context
		cancelFunc context.CancelFunc
	)
	ctx, cancelFunc = context.WithCancel(context.TODO())
	return &JobWorkInfo{
		Job:        job,
		ConText:    ctx,
		CancelFunc: cancelFunc,
	}
}

// slice 清空
func SliceClear(s *[]interface{}) {
	*s = (*s)[0:0]
}


// 获取当前的时间
func GetNowDate(format string) string {
	var (
		nTime time.Time
	)
	loc, _:= time.LoadLocation("Asia/Chongqing")
	nTime = time.Now().In(loc)
	return nTime.Format(format)
}

func CreateIndexByType(topic string, indexType string) string {
	var (
		dataFormat string
	)
	switch indexType {
	case "yyyy.MM.dd":
		dataFormat = "2006.01.02"
		break
	case "yyyy":
		dataFormat = "2006"
		break
	case "yyyy.MM":
		dataFormat = "2016.01"
	default:
		dataFormat = ""
	}
	if dataFormat != "" {
		return topic + "_" + GetNowDate(dataFormat)
	} else {
		return topic
	}
}

func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, conf.JobConf.JobSave)
}

func ExtractLockName(jobKey string) string {
	return strings.TrimPrefix(jobKey, conf.JobConf.JobLock)
}

