package logs

import (
	"encoding/json"
	"fmt"
	"log"
	"logserver/helper"
	"os"
	"time"
)

type File struct {
	logger *log.Logger
	file   *os.File
	init   bool
	fc *fileConf
	daily bool
}

type fileConf struct {
	Path       string `json:"path"`
	SaveName   string `json:"save_name"`
	LogFileExt string `json:"log_file_ext"`
}

func AddFile(daily bool, v ...interface{}) Logger {
	var  fc fileConf
	if len(v) == 0 {

	} else {
		if err := json.Unmarshal([]byte(v[0].(string)), &fc); err != nil {
			panic(err)
		}
	}
	return &File{
		daily:daily,
		fc: &fc,
	}
}

func (this *File) Init() {
	this.OpenFile()
	if this.daily {
		go this.dailyFile()
	}
}

func (this *File) OpenFile() {
	var (
		filePath string
		fc       *fileConf
	)
	if this.fc == nil {
		filePath = this.getLogFilePath() + "/" + this.getLogFileName()
	} else {
		fc = this.fc
		filePath = this.getLogFilePath(fc.Path) + "/" + this.getLogFileName(fc.SaveName, fc.LogFileExt)
	}
	if this.file != nil {
		this.file.Close()
	}
	this.file = this.openLogFile(filePath)
	this.logger = log.New(this.file, "", log.LstdFlags)
}

func (this *File) WriteMsg(pre string, v ...interface{}) {
	if this.init == false {
		this.Init()
	}
	this.logger.SetPrefix(pre)
	this.logger.Println(v...)
}

func (this *File) Destroy() {
	this.file.Close()
}

func (this *File) getLogFilePath(paths ...string) string {
	var path string
	if len(paths) == 0 {
		//path =
	} else {
		path = paths[0]
	}

	return fmt.Sprintf("%s/%s", helper.GetRootPath(), path)
}

// getLogFileName get the save name of the log file
func (this *File) getLogFileName(v ...string) string {
	var saveName, logFileExt, timily string
	if len(v) == 0 {

	} else {
		saveName = v[0]
		logFileExt = v[1]
	}
	if this.daily {
		timily = time.Now().Format("2006-01-02-15:04:05")
	}
	return fmt.Sprintf("%s%s.%s",
		saveName,
		timily,
		logFileExt,
	)
}

func (this *File) openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		this.mkDir()
	case os.IsPermission(err):
		log.Fatalf("Permission :%v", err)
	}

	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile :%v", err)
	}
	return handle
}

func (this *File) mkDir() {
	err := os.MkdirAll(this.getLogFilePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (this *File) dailyFile() {
	for {
		now := time.Now()
		//计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		this.OpenFile()
	}
}
