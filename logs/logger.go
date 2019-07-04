package logs

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

const (
	debug = iota
	info
	warning
	error
	fatal
)

type Logger interface {
	Init()
	WriteMsg(pre string, v ...interface{})
	Destroy()
}

type nameLogger struct {
	Logger
	name string
}

type MyLogger struct {
	mu          *sync.Mutex
	CallerDepth int
	outputs     []Logger
	level       int
	sign        chan os.Signal
	init        bool
	dataFormat  string
}

func NewMyLogger(level int) *MyLogger {
	mylogger := &MyLogger{
		sign:        make(chan os.Signal),
		mu:          new(sync.Mutex),
		level:       level,
		CallerDepth: 2,
		outputs:     make([]Logger, 0),
		init:        true,
		dataFormat:  "2006-01-02",
	}
	mylogger.Register(AddConsole())
	go mylogger.listenQuit()
	return mylogger
}

func (this *MyLogger) Register(Loggers ...Logger) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.outputs = append(this.outputs, Loggers...)
}

func (this *MyLogger) SetLevel(level int) {
	this.level = level
}

func (this *MyLogger) DEBUG(v ...interface{}) {
	if debug > this.level {
		return
	}
	this.wirteLog(debug, v...)
}

func (this *MyLogger) INFO(v ...interface{}) {
	if info > this.level {
		return
	}
	this.wirteLog(info, v...)
}

func (this *MyLogger) WARNING(v ...interface{}) {
	if warning > this.level {
		return
	}
	this.wirteLog(warning, v...)
}

func (this *MyLogger) ERROR(v ...interface{}) {
	if error > this.level {
		return
	}
	this.wirteLog(error, v...)
}

func (this *MyLogger) FATAL(v ...interface{}) {
	if fatal > this.level {
		return
	}
	this.wirteLog(fatal, v...)
}

func (this *MyLogger) wirteLog(level int, v ...interface{}) {
	var logPrefix string
	when := time.Now()
	_, file, line, ok := runtime.Caller(this.CallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d] %s ", levelFlags[level], filepath.Base(file), line, when.Format(this.dataFormat))
	} else {
		logPrefix = fmt.Sprintf("[%s] %s ", levelFlags[level], when.Format(this.dataFormat))
	}
	for _, output := range this.outputs {
		output.WriteMsg(logPrefix, v)
	}
}

func (this *MyLogger) listenQuit() {
	signal.Notify(this.sign, os.Interrupt)
	for {
		select {
		case <-this.sign:
			this.Destroy()
		}
	}
}

func (this *MyLogger) Destroy() {
	if this.init {
		for _, v := range this.outputs {
			v.Destroy()
		}
	}
	this.init = false
}

var myLogger = NewMyLogger(fatal)

func DEBUG(v ...interface{}) {
	myLogger.DEBUG(v...)
}

func INFO(v ...interface{}) {
	myLogger.INFO(v...)
}

func WARNING(v ...interface{}) {
	myLogger.WARNING(v...)
}

func ERROR(v ...interface{}) {
	myLogger.ERROR(v...)
}

func FATAL(v ...interface{}) {
	myLogger.FATAL(v...)
}

