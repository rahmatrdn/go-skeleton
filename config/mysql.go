package config

import (
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type Mysql struct {
	DB *gorm.DB
}

func NewMysql(env string, cfg *MysqlOption, dbLogger glogger.Interface) (*Mysql, error) {
	logLevel := glogger.Warn
	if env == "local" {
		logLevel = glogger.Info
	}

	db, err := gorm.Open(gmysql.Open(cfg.URI), &gorm.Config{
		Logger: dbLogger.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxOpenConns(cfg.Pool)
	return &Mysql{DB: db}, err
}
