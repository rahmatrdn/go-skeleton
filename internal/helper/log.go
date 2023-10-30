package helper

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rahmatrdn/go-skeleton/entity"
)

func WriteLogToFile(data string, channel string) error {
	f, err := os.OpenFile(channel,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer f.Close()

	_, err = f.WriteString(data)

	return err
}

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
	logFieldsJson, _ := json.Marshal(logFields)

	fmt.Println("=====================")
	fmt.Print(fmt.Sprintf("[%s] %s: ", DatetimeNowJakartaString(), status))

	logString := fmt.Sprintf("func_name='%s' error='%s' process='%s' message='%s' fields='%s'", logData.FuncName, logData.ErrorMessage, logData.Process, logData.Message, string(logFieldsJson))

	fmt.Println(logString)
	fmt.Println("=====================")

	ts := fmt.Sprintf("[%s] %s:", DatetimeNowJakartaString(), status)
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
