package config

import (
	"os"

	"github.com/rahmatrdn/go-skeleton/entity"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLog(env string) (*zap.Logger, error) {
	if env == entity.PRODUCTION_ENV {
		return NewProductionLogger()
	}

	return NewDevelopmentLogger()
}

// This configuration for Development env, the log will be written to the terminal!
func NewDevelopmentLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // ISO8601 time format
	config.EncoderConfig.TimeKey = "timestamp"
	return config.Build()
}

// This configuration for Production env, the log will be written to a File!
func NewProductionLogger() (*zap.Logger, error) {
	logDir := "storage/log"
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return nil, err
	}

	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // ISO8601 time format
	config.EncoderConfig.TimeKey = "timestamp"

	fileEncoder := zapcore.NewJSONEncoder(config.EncoderConfig)
	logFile, _ := os.OpenFile("storage/log/log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger, nil
}
