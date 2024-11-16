package config

import (
	"log"
	"os"
	"time"

	glogger "gorm.io/gorm/logger"
)

func NewGormLogMysqlConfig(cfg *MysqlOption) glogger.Interface {
	return glogger.New(
		log.New(
			os.Stdout,
			"\r\n",
			log.LstdFlags,
		),
		glogger.Config{
			SlowThreshold:             time.Duration(cfg.SlowThreshold) * time.Millisecond,
			LogLevel:                  glogger.Warn,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
		},
	)
}

func NewGormLogPostgreConfig(cfg *PostgreSqlOption) glogger.Interface {
	return glogger.New(
		log.New(
			os.Stdout,
			"\r\n",
			log.LstdFlags,
		),
		glogger.Config{
			SlowThreshold:             time.Duration(cfg.SlowThreshold) * time.Millisecond,
			LogLevel:                  glogger.Warn,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
		},
	)
}
