package config

import (
	"net/url"
	"strings"
	"time"

	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type Mysql struct {
	DB *gorm.DB
}

func NewMysql(env string, timezone string, cfg *MysqlOption, dbLogger glogger.Interface) (*Mysql, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}

	// Inject loc into DSN so the MySQL driver parses times in the correct timezone
	uri := cfg.URI
	if !strings.Contains(uri, "loc=") {
		sep := "&"
		if !strings.Contains(uri, "?") {
			sep = "?"
		}
		uri += sep + "loc=" + url.QueryEscape(timezone)
	}

	logLevel := glogger.Warn
	if env == "local" {
		logLevel = glogger.Info
	}

	db, err := gorm.Open(gmysql.Open(uri), &gorm.Config{
		Logger:  dbLogger.LogMode(logLevel),
		NowFunc: func() time.Time { return time.Now().In(loc) },
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxOpenConns(cfg.Pool)
	return &Mysql{DB: db}, err
}
