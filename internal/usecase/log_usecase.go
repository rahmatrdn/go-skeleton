package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/queue"
)

// LogUsecase is a usecase for writing log to Queue (Message Broker)
type Log struct {
	queue queue.Queue
}

func NewLogUsecase(
	queue queue.Queue,
) *Log {
	return &Log{queue}
}

type LogUsecase interface {
	Log(status entity.LogType, message string, funcName string, err error, logFields map[string]interface{}, processName string)
}

// Process writing log to file.
// Parameters :
//   - status: status of log (Check entity.LogType)
//   - message: message to descirbe the error (You can use it to indicate error dependencies/functions)
//   - funcName: source function that return error (Ex. walletUsecase.Create, etc.)
//   - err: error response from function
//   - logFields: additional data to track error (Ex. Indetifier ID, User ID, etc.)
//   - processName: name of process (optional, this can be use to track bug by process name) and make sure using Type Safety to write process name
func (w *Log) Log(status entity.LogType, message string, funcName string, err error, logFields map[string]interface{}, processName string) {
	logData := entity.Log{
		Process:      processName,
		FuncName:     funcName,
		Message:      message,
		ErrorMessage: err.Error(),
		Status:       status,
		LogFields:    logFields,
	}

	log, _ := json.Marshal(logData)

	ts := fmt.Sprintf("[%s] %s:", helper.NowStrUTC(), status)
	o := fmt.Sprint(string(log))

	payload, _ := helper.Serialize(logData)
	errQueue := w.queue.Publish(queue.ProcessSyncLog, payload, 1)

	// If error when publish to queue, write log to file
	if errQueue != nil {
		channel := "../../storage/log/general"
		channel = fmt.Sprintf("%s_%s.log", channel, helper.DateNowJakarta())
		f := fmt.Sprintf("%s %s \r\n", ts, o)
		helper.WriteLogToFile(f, channel)
		return
	}
}
