package config

import (
	"errors"
	"os"

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

// This configuration for Development env, the log will be written to the terminal!
func NewDevelopmentLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // ISO8601 time format
	config.EncoderConfig.TimeKey = "timestamp"
	config.DisableStacktrace = true           // Disable stack trace but keep caller info
	config.EncoderConfig.CallerKey = "caller" // Enable caller info for error location
	return config.Build()
}

func NewProductionLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // ISO8601 time format
	config.EncoderConfig.TimeKey = "timestamp"

	// Pastikan direktori ada
	if err := os.MkdirAll("storage/log", os.ModePerm); err != nil {
		return nil, err // Kembalikan error jika direktori gagal dibuat
	}

	// Coba buka file log
	logFile, err := os.OpenFile("storage/log/log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err // Kembalikan error jika gagal membuka file
	}

	fileEncoder := zapcore.NewJSONEncoder(config.EncoderConfig)
	writer := zapcore.AddSync(logFile)

	// Set level log default
	defaultLogLevel := zapcore.DebugLevel

	// Buat core zap tanpa stacktrace
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)

	// Inisialisasi logger tanpa stacktrace
	logger := zap.New(core, zap.AddCaller())

	if logger == nil {
		return nil, errors.New("failed to create logger")
	}

	return logger, nil
}
