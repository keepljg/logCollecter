package logs

import (
	"log"
	"os"
)

type ConSole struct {
	logger *log.Logger
	init bool
}

func AddConsole()Logger {
	return &ConSole{}
}


func (this *ConSole) Init () {
	this.logger = log.New(os.Stdout, "", 2)
}

func (this *ConSole) WriteMsg (pre string, v ...interface{}) {
	if this.init == false {
		this.Init()
	}
	this.logger.SetPrefix(pre)
	this.logger.Println(v...)
}

func (this *ConSole) Destroy() {

}