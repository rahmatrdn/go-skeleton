package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/rahmatrdn/go-skeleton/entity"
)

func WriteLogToFile(data string, channel string) error {
	f, err := os.OpenFile(channel,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	if _, err := f.WriteString(data); err != nil {
		log.Println(err)
	}

	return err
}

// func Log(fields logrus.Fields, logType string, logInfo string) {

// 	logger := logrus.New()
// 	logrus.SetFormatter(&logrus.JSONFormatter{})

// 	logger.SetOutput(os.Stdout)
// 	// logger.SetReportCaller(true)

// 	f, err := os.OpenFile("logrus.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err == nil {
// 		multi := io.MultiWriter(f, os.Stdout)
// 		logger.SetOutput(multi)
// 	} else {
// 		logger.SetOutput(os.Stdout)
// 	}

// 	switch logType {
// 	case "info":
// 		logger.WithFields(fields).Info(logInfo)
// 	case "warn":
// 		logger.WithFields(fields).Warn(logInfo)
// 	case "error":
// 		logger.WithFields(fields).Error(logInfo)
// 	default:
// 		logger.WithFields(fields).Info(logInfo)
// 	}
// }

func Log(status entity.LogType, message string, funcName string, err error, logFields entity.CaptureFields, processName string) {
	storage := entity.GeneralLogFilePath
	storage = fmt.Sprintf("%s_%s.log", storage, DateNowJakarta())

	logData := entity.Log{
		Process:      processName,
		FuncName:     funcName,
		Message:      message,
		ErrorMessage: err.Error(),
		Status:       status,
		LogFields:    logFields,
	}

	log, _ := json.Marshal(logData)

	fmt.Println("=====================")
	fmt.Println(fmt.Sprintf("[%s] %s:", DatetimeNowJakarta(), status), string(log))
	fmt.Println("=====================")

	ts := fmt.Sprintf("[%s] %s:", NowStrUTC(), status)
	o := fmt.Sprint(string(log))
	f := fmt.Sprintf("%s %s \r\n", ts, o)

	WriteLogToFile(f, storage)
}

func LogError(processName string, funcName string, err error, logFields entity.CaptureFields, message string) {
	Log(entity.LogError, message, funcName, err, logFields, processName)
}

func LogInfo(processName string, funcName string, logFields entity.CaptureFields, message string) {
	Log(entity.LogInfo, message, funcName, fmt.Errorf(""), logFields, processName)
}

func LogWarn(processName string, funcName string, err error, logFields entity.CaptureFields, message string) {
	Log(entity.LogWarning, message, funcName, err, logFields, processName)
}
