package usecase

import (
	"errors"
	"os"

	"github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/queue"
	"go.uber.org/zap"
)

// LogUsecase is a usecase for writing log to Queue (Message Broker)
type Log struct {
	queue     queue.Queue
	zapLogger *zap.Logger
}

func NewLogUsecase(
	queue queue.Queue,
	zapLogger *zap.Logger,
) *Log {
	return &Log{queue, zapLogger}
}

type LogUsecase interface {
	Log(status entity.LogType, message string, funcName string, err error, logFields map[string]string, processName string)
	Error(process string, funcName string, err error, logFields map[string]string)
	Info(message string, funcName string, logFields map[string]string, processName string)
}

// Process writing log to file.
// Parameters :
//   - status: status of log (Check entity.LogType)
//   - message: message to descirbe the error (You can use it to indicate error dependencies/functions)
//   - funcName: source function that return error (Ex. TodoListUsecase.Create, etc.)
//   - err: error response from function
//   - logFields: additional data to track error (Ex. Indetifier ID, User ID, etc.)
//   - processName: name of process (optional, this can be use to track bug by process name) and make sure using Type Safety to write process name
func (w *Log) Log(status entity.LogType, message string, funcName string, err error, logFields map[string]string, processName string) {
	logData := entity.Log{
		Process:      processName,
		FuncName:     funcName,
		Message:      message,
		ErrorMessage: err.Error(),
		Status:       status,
		LogFields:    logFields,
	}

	payload, _ := helper.Serialize(logData)
	errQueue := w.queue.Publish(queue.ProcessSyncLog, payload, 1)

	// Writing Log with Zap Logger
	logger := w.zapLogger.WithOptions(zap.AddCallerSkip(1))

	fields := []zap.Field{
		zap.String("process", processName),
		zap.String("funcName", funcName),
		zap.String("message", message),
		zap.String("errorMessage", err.Error()),
		zap.Any("logFields", logFields),
	}

	// If error when publish to queue, write log to file
	if errQueue != nil || (helper.GetAppEnv() != entity.PRODUCTION_ENV && os.Getenv("DEBUG_MODE") == "true") {
		switch status {
		case entity.LogError:
			logger.Error(message, fields...)
		case entity.LogInfo:
			logger.Info(message, fields...)
		}
		return
	}
}

func (w *Log) Error(process string, funcName string, err error, logFields map[string]string) {
	w.Log(entity.LogError, process, funcName, err, logFields, process)
}

func (w *Log) Info(message string, funcName string, logFields map[string]string, processName string) {
	w.Log(entity.LogInfo, message, funcName, errors.New(""), logFields, processName)
}
