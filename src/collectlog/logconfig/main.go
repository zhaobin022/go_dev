package logconfig

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"
)

func InitLog(logPath, level string) {
	config := make(map[string]interface{})
	config["filename"] = logPath

	fmt.Println("logpath : ", config["filename"])

	switch level {
	case "debug":
		config["level"] = logs.LevelDebug
	case "info":
		config["level"] = logs.LevelInfo
	case "error":
		config["level"] = logs.LevelError
	case "critical":
		config["level"] = logs.LevelCritical
	case "warning":
		config["level"] = logs.LevelWarning
	}
	fmt.Println("log level :", config["level"])
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("marshal failed, err:", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))
	logs.SetLogger("console")
	logs.SetLogFuncCall(true)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)

	// logs.Debug("this is a test, my name is %s", "stu01")
	// logs.Trace("this is a trace, my name is %s", "stu02")
	// logs.Warn("this is a warn, my name is %s", "stu03")
}
