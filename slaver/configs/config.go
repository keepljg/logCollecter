package configs

import (
	"fmt"
	"github.com/astaxie/beego/config"
)

var (
	AppConfig *Config
)

type Config struct {
	LogLevel    string
	LogPath     string
	JobSave     string
	JobLock     string
	ChanSize    int
	EsAddr      string
	KafkaAddr   string
	EtcdAddr    string
	EtcdKey     string
	EtcdTimeOut int
}

func InitConf(confType, filename string) (err error) {

	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		fmt.Println("new configs failed, err:", err)
		return
	}

	AppConfig = &Config{}
	AppConfig.LogLevel = conf.String("serverlogs::log_level")
	if len(AppConfig.LogLevel) == 0 {
		AppConfig.LogLevel = "debug"
	}

	AppConfig.LogPath = conf.String("serverlogs::log_path")
	if len(AppConfig.LogPath) == 0 {
		AppConfig.LogPath = "./serverlogs/logs.log"
	}

	AppConfig.JobSave = conf.String("job::job_save_dir")
	if len(AppConfig.JobSave) == 0 {
		AppConfig.JobSave = "/job/save/"
	}

	AppConfig.JobLock = conf.String("job::job_lock_dir")
	if len(AppConfig.JobLock) == 0 {
		AppConfig.JobLock = "/job/lock/"
	}

	AppConfig.ChanSize, err = conf.Int("collect::chan_size")
	if err != nil {
		AppConfig.ChanSize = 100
	}

	AppConfig.KafkaAddr = conf.String("kafka::server_addr")
	if len(AppConfig.KafkaAddr) == 0 {
		err = fmt.Errorf("invalid kafka addr")
		return
	}

	AppConfig.EsAddr = conf.String("es::addr")
	if len(AppConfig.EsAddr) == 0 {
		err = fmt.Errorf("invalid es addr")
		return
	}

	AppConfig.EtcdAddr = conf.String("etcd::addr")
	if len(AppConfig.EtcdAddr) == 0 {
		err = fmt.Errorf("invalid etcd addr")
		return
	}

	//AppConfig.EtcdKey = conf.String("etcd::configKey")
	//if len(AppConfig.EtcdKey) == 0 {
	//	err = fmt.Errorf("invalid etcd key")
	//	return
	//}

	AppConfig.EtcdTimeOut, err = conf.Int("etcd::timeout")
	if err != nil {
		AppConfig.EtcdTimeOut = 5
	}
	return
}
