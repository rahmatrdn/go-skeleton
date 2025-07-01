package helper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rahmatrdn/go-skeleton/config"
	"github.com/rahmatrdn/go-skeleton/entity"
	"go.uber.org/zap"
)

func WriteLogToFile(data string, channel string) error {
	dir := filepath.Dir(channel)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(channel,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	return err
}

// Function to Write Log
// If the app environment is set to production, the log will be written to a file.
// If the app environment is set to development, the log will be written to the terminal.
func Log(status entity.LogType, message string, funcName string, err error, logFields entity.CaptureFields, processName string) {
	logger, _ := config.NewZapLog(GetAppEnv())
	logger = logger.WithOptions(zap.AddCallerSkip(2))
	defer logger.Sync()

	fields := []zap.Field{
		zap.String("process", processName),
		zap.String("funcName", funcName),
		zap.String("message", message),
		zap.String("errorMessage", err.Error()),
		zap.Any("logFields", logFields),
	}

	switch status {
	case entity.LogError:
		logger.Error(message, fields...)
	case entity.LogInfo:
		logger.Info(message, fields...)
	case entity.LogDebug:
		logger.Debug(message, fields...)
	}

}

// Process writing log Error to file and console.
// Parameters :
//   - processName : name of process (optional, this can be use to track bug by process name) and make sure using Type Safety to write process name
//   - funcName : source function that return error (Ex. TodoListUsecase.Create, etc.)
//   - err : error response from function
//   - logFields : additional data to track error (Ex. Indetifier ID, User ID, etc.)
func LogError(process string, funcName string, err error, logFields entity.CaptureFields, message string) {
	Log(entity.LogError, process, funcName, err, logFields, process)
}

// Process writing log Info to file and console.
// Parameters :
//
//   - processName : name of process (optional, this can be use to track bug by process name) and make sure using Type Safety to write process name
//   - funcName : source function that return error (Ex. TodoListUsecase.Create, etc.)
//   - logFields : additional data to track error (Ex. Indetifier ID, User ID, etc.)
func LogInfo(processName string, funcName string, logFields entity.CaptureFields, message string) {
	Log(entity.LogInfo, message, funcName, fmt.Errorf(""), logFields, processName)
}

// Process writing log Warning to file and console.
// Parameters :
//   - processName : name of process (optional, this can be use to track bug by process name) and make sure using Type Safety to write process name
//   - funcName : source function that return error (Ex. TodoListUsecase.Create, etc.)
//   - err : error response from function
//   - logFields : additional data to track error (Ex. Indetifier ID, User ID, etc.)
func LogWarn(processName string, funcName string, err error, logFields entity.CaptureFields, message string) {
	Log(entity.LogWarning, message, funcName, err, logFields, processName)
}
