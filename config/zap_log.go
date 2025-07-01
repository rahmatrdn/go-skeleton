package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rahmatrdn/go-skeleton/entity"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLog(env string) (*zap.Logger, error) {
	if env == entity.PRODUCTION_ENV && os.Getenv("DEBUG_MODE") == "false" {
		return NewProductionLogger()
	}

	return NewDevelopmentLogger()
}

// NewDevelopmentLogger initializes and returns a zap.Logger configured for development use.
//
// It creates a structured JSON logger using zap's development configuration,
// including human-readable ISO8601 timestamps and caller information for easier debugging.
// Stack traces are disabled to keep the logs cleaner during local development,
// while the caller key is enabled to display the precise error location in your code.
func NewDevelopmentLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // ISO8601 time format
	config.EncoderConfig.TimeKey = "timestamp"
	config.DisableStacktrace = true           // Disable stack trace but keep caller info
	config.EncoderConfig.CallerKey = "caller" // Enable caller info for error location
	return config.Build()
}

// NewProductionLogger initializes and returns a zap.Logger configured for production use.
//
// It automatically creates a structured JSON logger with ISO8601 timestamp format,
// writing logs into daily log files located in a structured year/month folder hierarchy.
// For example, logs will be written into:
//
//	storage/log/2025/07/2025-07-01.log
//
// The function ensures that the necessary directories are created before writing,
// and logs are appended if the file already exists. The logger uses the DebugLevel
// by default, and includes caller information in each log entry for easier tracing.
func NewProductionLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // ISO8601 time format
	config.EncoderConfig.TimeKey = "timestamp"

	now := time.Now()
	yearMonth := now.Format("2006/01") // "2025/07"
	date := now.Format("2006-01-02")   // "2025-07-01"
	folderPath := filepath.Join("storage/log", yearMonth)
	logFileName := fmt.Sprintf("%s.log", date)
	fullLogPath := filepath.Join(folderPath, logFileName)

	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(fullLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	fileEncoder := zapcore.NewJSONEncoder(config.EncoderConfig)
	writer := zapcore.AddSync(logFile)

	defaultLogLevel := zapcore.DebugLevel

	core := zapcore.NewCore(fileEncoder, writer, defaultLogLevel)

	logger := zap.New(core, zap.AddCaller())

	if logger == nil {
		return nil, errors.New("failed to create logger")
	}

	return logger, nil
}
