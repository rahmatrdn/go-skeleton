package helper

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rahmatrdn/go-skeleton/entity"
)

func WriteLogToFile(data string, channel string) error {
	f, _ := os.OpenFile(channel,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer f.Close()

	_, err := f.WriteString(data)

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

	fmt.Println("=====================")
	fmt.Printf("[%s] %s: ", DatetimeNowJakartaString(), status)

	logString := fmt.Sprintf("func_name='%s' error='%s' process='%s' message='%s' fields='%s'", logData.FuncName, logData.ErrorMessage, logData.Process, logData.Message, logFields)

	fmt.Println(logString)
	fmt.Println("=====================")

	ts := fmt.Sprintf("[%s] %s:", DatetimeNowJakartaString(), status)
	o := fmt.Sprint(string(log))
	f := fmt.Sprintf("%s %s \r\n", ts, o)

	WriteLogToFile(f, storage)
}

// Process writing log Error to file and console.
// Parameters :
//   - processName : name of process (optional, this can be use to track bug by process name) and make sure using Type Safety to write process name
//   - funcName : source function that return error (Ex. walletUsecase.Create, etc.)
//   - err : error response from function
//   - logFields : additional data to track error (Ex. Indetifier ID, User ID, etc.)
func LogError(processName string, funcName string, err error, logFields entity.CaptureFields, message string) {
	Log(entity.LogError, message, funcName, err, logFields, processName)
}

// Process writing log Info to file and console.
// Parameters :
//
//   - processName : name of process (optional, this can be use to track bug by process name) and make sure using Type Safety to write process name
//   - funcName : source function that return error (Ex. walletUsecase.Create, etc.)
//   - logFields : additional data to track error (Ex. Indetifier ID, User ID, etc.)
func LogInfo(processName string, funcName string, logFields entity.CaptureFields, message string) {
	Log(entity.LogInfo, message, funcName, fmt.Errorf(""), logFields, processName)
}

// Process writing log Warning to file and console.
// Parameters :
//   - processName : name of process (optional, this can be use to track bug by process name) and make sure using Type Safety to write process name
//   - funcName : source function that return error (Ex. walletUsecase.Create, etc.)
//   - err : error response from function
//   - logFields : additional data to track error (Ex. Indetifier ID, User ID, etc.)
func LogWarn(processName string, funcName string, err error, logFields entity.CaptureFields, message string) {
	Log(entity.LogWarning, message, funcName, err, logFields, processName)
}
