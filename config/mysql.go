package config

import (
	"fmt"

	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type Mysql struct {
	DB *gorm.DB
}

func NewMysql(env string, cfg *MysqlOption, dbLogger glogger.Interface) (*Mysql, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DatabaseName,
		// url.QueryEscape(cfg.TimeZone),
	)

	logLevel := glogger.Warn
	if env == "local" {
		logLevel = glogger.Info
	}

	db, err := gorm.Open(gmysql.Open(dsn), &gorm.Config{
		Logger: dbLogger.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxOpenConns(cfg.Pool)
	return &Mysql{DB: db}, err
}
