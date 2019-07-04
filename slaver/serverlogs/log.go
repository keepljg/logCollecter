package serverlogs

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"logserver/slaver/conf"
)

func convertLogLevel(level string) int {

	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}

	return logs.LevelDebug
}

func InitLogger() (err error) {

	config := make(map[string]interface{})
	config["filename"] = conf.AppConfig.LogPath
	config["level"] = convertLogLevel(conf.AppConfig.LogLevel)

	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("initLogger failed, marshal err:", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}
